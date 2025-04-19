package imple

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ntquang/ecommerce/global"
	consts "github.com/ntquang/ecommerce/internal/const"
	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
	utlis "github.com/ntquang/ecommerce/internal/utils"
	"github.com/ntquang/ecommerce/internal/utils/auth"
	"github.com/ntquang/ecommerce/internal/utils/crypto"
	"github.com/ntquang/ecommerce/internal/utils/random"
	"github.com/ntquang/ecommerce/internal/utils/sendto"
	"github.com/ntquang/ecommerce/response"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type sUserLogin struct {
	// implement the IUserLogin interface here
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}

// implement the IUserLogin interface here

// --- START TWO FACTOR AUTHENTICATION ---
func (s *sUserLogin) IsTwoFactorEnabled(ctx context.Context, userId int) (resultCode int, rs bool, err error) {
	return 200, true, nil
}
func (s *sUserLogin) SetupTwoFactorAuth(ctx context.Context, in *model.SetupTwoFactorAuthInput) (resultCode int, err error) {
	// 1. check user enabled
	isTwoFactorAuth, err := s.r.IsTwoFactorEnabled(ctx, int32(in.UserId))
	if err != nil {
		return response.ErrTwoFactorAuthSetUpFailed, err
	}
	if isTwoFactorAuth > 0 {
		return response.ErrTwoFactorAuthSetUpFailed, fmt.Errorf("Two-factor authenticaion already enabled")
	}
	// 2. create new type Authen
	err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
		UserID:            int32(in.UserId),
		TwoFactorAuthType: database.TwoFactorAuthTypeEnumEMAIL,
		TwoFactorEmail:    pgtype.Text{String: in.TwoFactorEmail, Valid: true},
	})
	if err != nil {
		return response.ErrTwoFactorAuthSetUpFailed, err
	}
	// 3. send otp to in.TwoFactorEmail
	keyUserTwoFactor := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))
	fmt.Println("keyUserTwoFactor::", keyUserTwoFactor)
	go global.Redis.Set(ctx, keyUserTwoFactor, "123456", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()

	return response.ErrCodeSuccess, nil
}
func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *model.TwoFactorVerifycationInput) (resultCode int, err error) {
	// 1. check two factor authentication enable
	isTwoFactorAuth, err := s.r.IsTwoFactorEnabled(ctx, int32(in.UserId))
	if err != nil {
		return response.ErrTwoFactorAuthVerifyFailed, fmt.Errorf("Two-factor authenticaion verify error")
	}
	if isTwoFactorAuth > 0 {
		return response.ErrTwoFactorAuthVerifyFailed, fmt.Errorf("Two-factor authenticaion is enabled")
	}

	// 2. check Otp redis avalible
	keyUserTwoFactor := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))
	otpVerifyAuth, err := global.Redis.Get(ctx, keyUserTwoFactor).Result()
	if err != nil {
		return response.ErrTwoFactorAuthVerifyFailed, fmt.Errorf("Two-factor authenticaion get key in redis error, %s", err)
	} else if err == redis.Nil {
		return response.ErrTwoFactorAuthVerifyFailed, fmt.Errorf("Two-factor authenticaion get key in redis empty")
	}
	fmt.Println("otpVerifyAuth", otpVerifyAuth)
	// 3. check otp match
	if otpVerifyAuth != in.TwoFactorCode {
		return response.ErrTwoFactorAuthVerifyFailed, fmt.Errorf("Two-factor authenticaion code not match")
	}

	// 4. update status
	err = s.r.UpdateTwoFactorStatusVerification(ctx, database.UpdateTwoFactorStatusVerificationParams{
		UserID:            int32(in.UserId),
		TwoFactorAuthType: database.TwoFactorAuthTypeEnumEMAIL,
	})
	if err != nil {
		return response.ErrTwoFactorAuthVerifyFailed, err
	}

	// 5. remove otp
	_, err = global.Redis.Del(ctx, keyUserTwoFactor).Result()
	if err != nil {
		return response.ErrTwoFactorAuthVerifyFailed, err
	}
	return 200, nil
}

// --- END TWO FACTOR AUTHENTICATION ---

func (s *sUserLogin) Login(ctx context.Context, in *model.LoginInput) (resultCode int, out model.LoginOutput, err error) {
	///1. check user in model user base
	userBase, err := s.r.GetOneUserInfo(ctx, in.Account)
	if err != nil {
		return response.ErrCodeUserNotRegister, out, err
	}
	// 2. check match password
	if !crypto.MatchingPassword(in.Password, userBase.UserSalt, userBase.UserPassword) {
		return response.ErrCodeUserNotRegister, out, err
	}
	// 3. check two-factor authentication
	isTwofactorEnable, err := s.r.IsTwoFactorEnabled(ctx, userBase.UserID)
	if err != nil {
		return response.ErrCodeAuthenError, out, err
	}
	if isTwofactorEnable > 0 {
		// 3.1 create key and set in rdb
		keyUserTwoFactor := crypto.GetHash("2fa:" + strconv.Itoa(int(userBase.UserID)))
		err := global.Redis.SetEx(ctx, keyUserTwoFactor, "11111", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
		if err != nil {
			return response.ErrCodeRedisSetFailed, out, fmt.Errorf("set otp to redis error")
		}
		// 3.2 get email 2fa
		infoUserTwoFactor, err := s.r.GetTwoFactorMethodByIDandType(ctx, database.GetTwoFactorMethodByIDandTypeParams{
			UserID:            userBase.UserID,
			TwoFactorAuthType: database.TwoFactorAuthTypeEnumEMAIL,
		})
		if err != nil {
			return response.ErrCodeAuthenError, out, err
		}

		// 3.3 send email
		log.Println("send OTP to Email::", infoUserTwoFactor.TwoFactorEmail)
		go sendto.SendTemplateEmailOtp(
			[]string{infoUserTwoFactor.TwoFactorEmail.String},
			consts.EMAIL_SEND,
			"otp-auth.html",
			map[string]interface{}{
				"otp":  "11111",
				"name": "Quang",
			},
		)
		out.Message = "Send OTP 2FA To Email, pls get otp by email..."
		return response.ErrCodeSuccess, out, err
	}

	//4. update password in user base (time, ip, ...)
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserLoginIp:  pgtype.Text{String: "172.0.0.1", Valid: true},
		UserAccount:  in.Account,
		UserPassword: in.Password,
	})
	//5. create uuid
	newUUID := utlis.GenerateCliTokenUUID(int(userBase.UserID))
	fmt.Println("uuid::", newUUID)
	//6. get user info
	userInfo, err := s.r.GetUser(ctx, int64(userBase.UserID))
	if err != nil {
		global.Logger.Info("error::", zap.Error(err))
		return response.ErrCodeUserNotRegister, out, err
	}
	fmt.Println("userInfo::", userInfo)
	// 7. convert user info to json and save userinfo in redis as uuid
	UserInfoJson, err := json.Marshal(userInfo)
	err = global.Redis.Set(ctx, newUUID, UserInfoJson, time.Duration(consts.TIME_2FA_LOGIN)*time.Minute).Err()
	//8. create token
	out.Token, err = auth.CreateToken(newUUID)
	if err != nil {
		return
	}
	return 200, out, err
}

func (s *sUserLogin) Register(ctx context.Context, in *model.RegisterInput) (resultCode int, err error) {
	fmt.Printf("VerifyKey: %s\n", in.VerifyKey)
	fmt.Printf("VerifyType: %d\n", in.VerifyType)
	// 0. hash VerifyKey
	hashKey := crypto.GetHash(in.VerifyKey)
	fmt.Printf("hashKey: %s\n", hashKey)
	// 1 check user exists
	userFound, err := s.r.CheckUserBaseExists(ctx, in.VerifyKey)
	if err != nil {
		return response.ErrCodeUserHasExists, err
	}

	fmt.Println("USER::", userFound)

	if userFound > 0 {
		return response.ErrCodeUserHasExists, fmt.Errorf("user has already registered")
	}

	//3 create OTP
	userKey := utlis.GetUserKey(hashKey)
	otpFound, err := global.Redis.Get(ctx, userKey).Result()

	switch {
	case err == redis.Nil:
		fmt.Println("Key doesn't exists!")
	case err != nil:
		fmt.Println("get failed ::", err)
		return response.ErrInvalidOtp, err
	case otpFound != "":
		return response.ErrCodeOtpNotExists, fmt.Errorf("")
	}

	//4 generate OTP
	otpNew := random.GenerateSixDigitOtp()
	if in.VerifyPurpose == "TEST_USER" {
		otpNew = 123456
	}
	fmt.Printf("otpNew: %s\n", otpNew)

	//5. save OTP in redos
	err = global.Redis.Set(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrInvalidOtp, err
	}

	//6 method register
	switch in.VerifyType {
	case consts.EMAIL:
		// 6.1 sent otp use go
		err = sendto.SendTextEmailOtp([]string{in.VerifyKey}, consts.EMAIL_SEND, strconv.Itoa(otpNew))
		if err != nil {
			return response.ErrSendOtp, err
		}
		// 6.2 send otp via kafka any service
		// body := make(map[string]interface{})
		// // khởi tạo và cấp phát bộ nhớ như new nhưng khác là chỉ cấp phát bộ nhớ nhưng không khởi tạo giá trị.
		// body["otp"] = otpNew
		// body["email"] = consts.EMAIL_SEND

		// bodyRequest, _ := json.Marshal(body)

		// message := kafka.Message{
		// 	Key:   []byte("oth-auth"),
		// 	Value: []byte(bodyRequest),
		// 	Time:  time.Now(),
		// }

		// err := global.KafkaProducer.WriteMessages(context.Background(), message)
		// if err != nil {
		// 	return response.ErrKafkaSendMessageFailed, err
		// }
		// 7. save OTP to Postgresql
		result, err := s.r.InsertOtpVerify(ctx, database.InsertOtpVerifyParams{
			VerifyOtp:     strconv.Itoa(otpNew),
			VerifyKey:     in.VerifyKey,
			VerifyType:    pgtype.Int4{Int32: 1, Valid: true},
			VerifyKeyHash: hashKey,
		})
		if err != nil {
			return response.ErrSendOtp, err
		}
		//8. getlastId when have function work need last ID
		// lastIdVerifyUser, err := result.

		fmt.Println("lastId ::", result)
		return response.ErrCodeSuccess, nil
	case consts.MOBILE:
		return response.ErrCodeSuccess, nil
	}
	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyOtp(ctx context.Context, in *model.VerifyInput) (out model.VerifyOTPOutput, err error) {
	// 1 hash email
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))

	// 2 get otp
	userKey := utlis.GetUserKey(hashKey)
	otpFound, err := global.Redis.Get(ctx, userKey).Result()
	if err != nil {
		return out, err
	}

	if in.VerifyCode != otpFound {
		// if failed count as three throw error
		// ... login
		return out, fmt.Errorf("Otp not match")
	}

	infoOTP, err := s.r.GetInfoOtp(ctx, hashKey)
	if err != nil {
		return out, err
	}

	// 3. update status
	err = s.r.UpdateUserValificationStatus(ctx, hashKey)
	if err != nil {
		return out, fmt.Errorf("Update failed!!")
	}

	out.Token = infoOTP.VerifyKeyHash
	out.Message = "Verify OTP successfull"
	return out, err
}

func (s *sUserLogin) UpdatePaswordRegister(ctx context.Context, token string, password string) (userId int, err error) {
	// 1. check token exists in dbs
	InfoOTP, err := s.r.GetInfoOtp(ctx, token)
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}
	// 2. check OTP already verify ?
	if InfoOTP.IsVerified.Int32 == 0 {
		return response.ErrCodeOtpNotVerify, err
	}

	//update userBase table
	userBase := database.AddUserBaseParams{}
	userBase.UserAccount = InfoOTP.VerifyKey
	userSalt, err := crypto.GenerateSalt(16)
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}
	userBase.UserSalt = userSalt
	userBase.UserPassword = crypto.HashPassword(password, userSalt)
	//add user Base to table
	newUserBase, err := s.r.AddUserBase(ctx, userBase)
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}
	fmt.Println("user Id::", newUserBase)
	//add user_id to unserInfo
	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:       int64(newUserBase),
		UserAccount:  InfoOTP.VerifyKey,
		UserNickname: pgtype.Text{String: "", Valid: true},
		UserAvatar:   pgtype.Text{String: "", Valid: true},
		UserState:    1,
		UserMobile:   pgtype.Text{String: "", Valid: true},
		UserGender:   pgtype.Int2{Int16: 0, Valid: true},
		UserBirthday: pgtype.Date{Time: time.Time{}, Valid: false},
		UserEmail:    pgtype.Text{String: InfoOTP.VerifyKey, Valid: true},
	})
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}
	return int(newUserInfo), nil
}
