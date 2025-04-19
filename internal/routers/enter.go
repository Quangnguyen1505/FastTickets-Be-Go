package routers

import (
	contactmessage "github.com/ntquang/ecommerce/internal/routers/contact_message"
	"github.com/ntquang/ecommerce/internal/routers/event"
	"github.com/ntquang/ecommerce/internal/routers/manage"
	menuFunction "github.com/ntquang/ecommerce/internal/routers/menu_function"
	"github.com/ntquang/ecommerce/internal/routers/oauth2"
	"github.com/ntquang/ecommerce/internal/routers/user"
)

type RouterGroup struct {
	Manage              manage.ManageRouterGroup
	User                user.UserRouterGroup
	Oauth2              oauth2.Oauth2RouterGroup
	Event               event.EventGroup
	MenuFunction        menuFunction.MenuFunctionGroup
	ContactMessageGroup contactmessage.ContactMessageGroup
}

var RouterGroupApp = &RouterGroup{}
