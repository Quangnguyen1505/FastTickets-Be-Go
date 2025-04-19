package user

import "github.com/gin-gonic/gin"

type ProductRouter struct{}

func (pr *ProductRouter) InitProductRouter(Router *gin.RouterGroup) {

	productRouterPublic := Router.Group("/product")
	{
		productRouterPublic.GET("/search")
		productRouterPublic.GET("/getCurrent")
	}

	productRouterPrivate := Router.Group("/product")
	// productRouterPrivate.Use(Limiter())
	// productRouterPrivate.Use(Authen())
	// productRouterPrivate.Use(Permission())
	{
		productRouterPrivate.GET("/create")
	}
}
