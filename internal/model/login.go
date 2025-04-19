package model

type RegisterInput struct {
	VerifyKey     string `json:"verify_key"`
	VerifyType    int    `json:"verify_type"`
	VerifyPurpose string `json:"verify_purpose"`
}

type VerifyInput struct {
	VerifyCode string `json:"verify_code"`
	VerifyKey  string `json:"verify_key"`
}

type VerifyOTPOutput struct {
	Token   string `json:"token"`
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

type UpdatePasswordInput struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type LoginInput struct {
	Account  string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	UserId  string `json:"user_id"`
}

type SetupTwoFactorAuthInput struct {
	UserId            uint32 `json:"user_id"`
	TwoFactorAuthType string `json:"two_factor_auth_type"`
	TwoFactorEmail    string `json:"two_factor_email"`
}

type TwoFactorVerifycationInput struct {
	UserId        uint32 `json:"user_id"`
	TwoFactorCode string `json:"two_factor_code"`
}
