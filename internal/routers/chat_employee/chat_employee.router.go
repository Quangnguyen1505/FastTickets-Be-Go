package chatemployee

import (
	"github.com/gin-gonic/gin"
	chatemployee "github.com/ntquang/ecommerce/internal/controller/chat_employee"
	"github.com/ntquang/ecommerce/internal/middlewares"
)

type ChatWithEmployeeRouter struct{}

func (cER *ChatWithEmployeeRouter) InitChatWithEmployee(Router *gin.RouterGroup) {
	publicChatWithEmployeeRouter := Router.Group("/chat-with-employee")
	publicChatWithEmployeeRouter.Use(middlewares.Authentication())
	{
		publicChatWithEmployeeRouter.POST("/init-session", chatemployee.ChatEmployee.InitSessionId)
		publicChatWithEmployeeRouter.GET("/history/:sessionId", chatemployee.ChatEmployee.GetHistoryChat)
		publicChatWithEmployeeRouter.POST("/send-message", chatemployee.ChatEmployee.SendMessage)
	}

	privateChatWithEmployeeRouter := Router.Group("/chat-with-employee")
	privateChatWithEmployeeRouter.Use(middlewares.Authentication())
	privateChatWithEmployeeRouter.Use(middlewares.CheckPermission())
	{
		privateChatWithEmployeeRouter.GET("/all-session", chatemployee.ChatEmployee.GetAllSessionId)
		privateChatWithEmployeeRouter.POST("/join-session", chatemployee.ChatEmployee.JoinChatSession)
		privateChatWithEmployeeRouter.GET("/close-session/:sessionId", chatemployee.ChatEmployee.CloseChatSession)
	}
}
