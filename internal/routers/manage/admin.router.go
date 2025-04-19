package manage

import (
	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/internal/controller/account"
)

type AdminRouter struct{}

func (pr *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {

	adminRouterPublic := Router.Group("/admin")
	{
		adminRouterPublic.GET("/login")
		adminRouterPublic.GET("/remove", account.UserAdmin.RemoveUser)
	}

	adminRouterPrivate := Router.Group("/admin")
	// adminRouterPrivate.Use(Limiter())
	// adminRouterPrivate.Use(Authen())
	// adminRouterPrivate.Use(Permission())
	{
		adminRouterPrivate.GET("/active_user")
	}
}
