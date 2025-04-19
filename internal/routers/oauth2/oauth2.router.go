package oauth2

import (
	"github.com/gin-gonic/gin"
	oauth2Controller "github.com/ntquang/ecommerce/internal/controller/oauth2"
)

type Oauth2Router struct{}

func (oauth2 *Oauth2Router) InitOauth2Router(Router *gin.RouterGroup) {
	oauth2RouterPublic := Router.Group("/auth")
	{
		oauth2RouterPublic.GET("/:provider", oauth2Controller.Oauth2.LoginWithProvider)
		oauth2RouterPublic.GET("/:provider/callback", oauth2Controller.Oauth2.CallbackHandler)
	}
}
