package routers

import (
	chatemployee "github.com/ntquang/ecommerce/internal/routers/chat_employee"
	contactmessage "github.com/ntquang/ecommerce/internal/routers/contact_message"
	"github.com/ntquang/ecommerce/internal/routers/event"
	menuFunction "github.com/ntquang/ecommerce/internal/routers/menu_function"
	"github.com/ntquang/ecommerce/internal/routers/oauth2"
)

type RouterGroup struct {
	Oauth2              oauth2.Oauth2RouterGroup
	Event               event.EventGroup
	MenuFunction        menuFunction.MenuFunctionGroup
	ContactMessageGroup contactmessage.ContactMessageGroup
	ChatEmployeeGroup   chatemployee.ChatEmployeeGroup
}

var RouterGroupApp = &RouterGroup{}
