package contactmessage

import (
	"github.com/gin-gonic/gin"
	contactmessage "github.com/ntquang/ecommerce/internal/controller/contact_message"
	"github.com/ntquang/ecommerce/internal/middlewares"
)

type ContactMessageRouter struct{}

func (cMR *ContactMessageRouter) InitContactMessage(Router *gin.RouterGroup) {
	publicContactMessageRouter := Router.Group("/contact-messages")
	{
		publicContactMessageRouter.POST("", contactmessage.ContactMessage.NewContactMessage)
		publicContactMessageRouter.GET("/:id", contactmessage.ContactMessage.GetContactMessageById)
	}

	privateContactMessageRouter := Router.Group("/contact-messages")
	privateContactMessageRouter.Use(middlewares.Authentication())
	{
		privateContactMessageRouter.GET("", contactmessage.ContactMessage.GetAllContactMessageStatus)
		privateContactMessageRouter.PUT("/:id", contactmessage.ContactMessage.EditContactMessage)
		privateContactMessageRouter.DELETE("/:id", contactmessage.ContactMessage.DeleteContactMessage)
	}
}
