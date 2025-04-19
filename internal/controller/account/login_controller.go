package account

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/internal/services"
	"github.com/ntquang/ecommerce/response"
	"go.uber.org/zap"
)

var Login = new(cUserLogin)

type cUserLogin struct{}

// User Login Documentation
// @Summary      User Login
// @Description  When user is login save dbs
// @Tags         account manager
// @Accept       json
// @Produce      json
// @Param        payload body model.LoginInput true "payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /users/login [post]
func (c *cUserLogin) Login(ctx *gin.Context) {
	//implement login for register
	var params model.LoginInput
	if err := ctx.ShouldBindJSON(&params); err != nil { // convert JSON to Struct
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "params invalid", fmt.Errorf(err.Error()))
	}
	statusCode, metadata, err := services.UserLogin().Login(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "Error login user", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, metadata, "login successfully!")
}

// User Registrasion Documentation
// @Summary      User Registrasion
// @Description  When user is registered send otp to email
// @Tags         account manager
// @Accept       json
// @Produce      json
// @Param        payload body model.RegisterInput true "payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /users/register [post]
func (c *cUserLogin) Register(ctx *gin.Context) {
	//implement login for register
	var params model.RegisterInput
	fmt.Println("so 1::", params)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error registering user OTP", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}
	fmt.Println("so 2::", params)
	statusCode, err := services.UserLogin().Register(ctx, &params)
	if err != nil {
		global.Logger.Error("Error registering user OTP", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, nil, "register successfully!")
}

// User Verify OTP Documentation
// @Summary      User Verify OTP
// @Description  When user is register after verify otp
// @Tags         account manager
// @Accept       json
// @Produce      json
// @Param        payload body model.VerifyInput true "payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /users/verifyOTP [post]
func (c *cUserLogin) VerifyOTP(ctx *gin.Context) {
	var params model.VerifyInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
	}
	result, err := services.UserLogin().VerifyOtp(ctx, &params)
	if err != nil {
		global.Logger.Error("Error verify user OTP", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, result, "verify otp successfully!")
}

// User Verify OTP Documentation
// @Summary      User Update Password
// @Description  When user is VerifyOTP ok after Update Password
// @Tags         account manager
// @Accept       json
// @Produce      json
// @Param        payload body model.UpdatePasswordInput true "payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /users/updatePass [post]
func (c *cUserLogin) UpdatePasswordRegister(ctx *gin.Context) {
	var params model.UpdatePasswordInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
	}
	userId, err := services.UserLogin().UpdatePaswordRegister(ctx, params.Token, params.Password)
	if err != nil {
		global.Logger.Error("Error update password", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, userId, "update password successfully")
}
