package menufunction

import (
	"github.com/gin-gonic/gin"
	MenuFunc "github.com/ntquang/ecommerce/internal/controller/menu_function"
	"github.com/ntquang/ecommerce/internal/middlewares"
)

type MenuFunctionRouter struct{}

func (menuFunc *MenuFunctionRouter) InitMenufunctionRouter(Router *gin.RouterGroup) {
	publicMenuFuntion := Router.Group("/menu-functions")
	{
		publicMenuFuntion.GET("/active", MenuFunc.MenuFunc.GetAllMenuFunctionsActive)
		publicMenuFuntion.GET("/:id", MenuFunc.MenuFunc.GetMenuFunctionsById)
	}

	privateMenuFuntion := Router.Group("/menu-functions")
	privateMenuFuntion.Use(middlewares.Authentication())
	privateMenuFuntion.Use(middlewares.CheckPermission())
	{
		privateMenuFuntion.GET("", MenuFunc.MenuFunc.GetAllMenuFunctions)
		privateMenuFuntion.POST("", MenuFunc.MenuFunc.NewMenuFunctions)
		privateMenuFuntion.PUT("/:id", MenuFunc.MenuFunc.EditMenuFunctionsById)
		privateMenuFuntion.DELETE("/:id", MenuFunc.MenuFunc.DeleteMenuFunctionsById)
	}
}
