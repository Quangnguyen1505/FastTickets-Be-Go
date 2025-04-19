package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/internal/controller/account"
	"github.com/ntquang/ecommerce/internal/middlewares"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {

	//middleware rate limit global
	Router.Use(middlewares.NewRateLimiterV2().GlobalRateLimiterV2())

	userRouterPublic := Router.Group("/users")
	userRouterPublic.Use(middlewares.NewRateLimiterV2().PublicAPIRateLimiterV2())
	{
		userRouterPublic.GET("/info-public")
		// userRouterPublic.POST("/register", userController.Register)
		userRouterPublic.POST("/register", account.Login.Register)
		userRouterPublic.POST("/verifyOTP", account.Login.VerifyOTP)
		userRouterPublic.POST("/login", account.Login.Login)
		userRouterPublic.POST("/updatePass", account.Login.UpdatePasswordRegister)
	}

	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(middlewares.Authentication())
	userRouterPrivate.Use(middlewares.NewRateLimiterV2().UserPrivateAPIRateLimiterV2())
	{
		userRouterPrivate.GET("/info")
		userRouterPrivate.POST("/two-factor/setup", account.User2fa.SetupTwoFactorAuth)
		userRouterPrivate.POST("/two-factor/verify", account.User2fa.VerifyTwoFactorAuth)
	}
}
