package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/middlewares"
	"github.com/ntquang/ecommerce/internal/routers"
	"github.com/ntquang/ecommerce/internal/websocket"
)

func Initrouter() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Logger(), gin.Recovery())
	}

	r.GET("/ws", func(c *gin.Context) {
		websocket.ChatSocketHandler(c.Writer, c.Request)
	})

	oauth2Router := routers.RouterGroupApp.Oauth2
	eventRouter := routers.RouterGroupApp.Event
	menuFunctionRouter := routers.RouterGroupApp.MenuFunction
	contactMessageRouter := routers.RouterGroupApp.ContactMessageGroup
	chatEmployee := routers.RouterGroupApp.ChatEmployeeGroup

	r.Use(middlewares.CORSMiddleware())
	MainGroup := r.Group("/v1/2024")
	{
		MainGroup.GET("/checkStatus") //tracking monitor
	}
	{
		oauth2Router.InitOauth2Router(MainGroup)
	}
	{
		eventRouter.InitEventRouter(MainGroup)
	}
	{
		menuFunctionRouter.InitMenufunctionRouter(MainGroup)
	}
	{
		contactMessageRouter.InitContactMessage(MainGroup)
	}
	{
		chatEmployee.InitChatWithEmployee(MainGroup)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route Not Found"})
	})

	return r
}
