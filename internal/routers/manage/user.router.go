package manage

import "github.com/gin-gonic/gin"

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {

	// userRouterPublic := Router.Group("/user")
	// {
	// 	userRouterPublic.GET("/register")
	// 	userRouterPublic.GET("/otp")
	// }

	userRouterPrivate := Router.Group("/admin/user")
	// userRouterPrivate.Use(Limiter())
	// userRouterPrivate.Use(Authen())
	// userRouterPrivate.Use(Permission())
	{
		userRouterPrivate.GET("/active_user")
	}
}
