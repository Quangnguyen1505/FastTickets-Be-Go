package services

import (
	"context"

	"github.com/ntquang/ecommerce/internal/model"
)

type (
	IUserLogin interface {
		Login(ctx context.Context, in *model.LoginInput) (resultCode int, out model.LoginOutput, err error)
		Register(ctx context.Context, in *model.RegisterInput) (resultCode int, err error)
		VerifyOtp(ctx context.Context, in *model.VerifyInput) (out model.VerifyOTPOutput, err error)
		UpdatePaswordRegister(ctx context.Context, token string, password string) (userId int, err error)

		//TwoFactorAuthentication
		IsTwoFactorEnabled(ctx context.Context, userId int) (resultCode int, rs bool, err error)
		//Setup authentication
		SetupTwoFactorAuth(ctx context.Context, in *model.SetupTwoFactorAuthInput) (resultCode int, err error)
		//verify two factor authentication
		VerifyTwoFactorAuth(ctx context.Context, in *model.TwoFactorVerifycationInput) (resultCode int, err error)
	}

	IUserInfo interface {
		GetInfoUserById(ctx context.Context) error
		GetAllUser(ctx context.Context) error
	}

	IUserAdmin interface {
		RemoveUser(ctx context.Context) (result string)
		FindOneUser(ctx context.Context) error
	}
)

var (
	localUserAdmin IUserAdmin
	localUserInfo  IUserInfo
	localUserLogin IUserLogin
)

func UserAdmin() IUserAdmin {
	if localUserAdmin == nil {
		panic("implement localUserAdmin not found for interface IUserAdmin")
	}

	return localUserAdmin
}

func InitUserAdmin(i IUserAdmin) {
	localUserAdmin = i
}

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("implement localUserInfo not found for interface IUserInfo")
	}

	return localUserInfo
}

func InitUserInfo(i IUserInfo) {
	localUserInfo = i
}

func UserLogin() IUserLogin {
	if localUserLogin == nil {
		panic("implement localUserLogin not found for interface IUserLogin")
	}

	return localUserLogin
}

func InitUserLogin(i IUserLogin) {
	localUserLogin = i
}
