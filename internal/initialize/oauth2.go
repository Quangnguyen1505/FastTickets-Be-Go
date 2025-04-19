package initialize

import (
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/ntquang/ecommerce/global"
)

func InitOauth2() {
	clientID := global.Config.Oauth2Google.CLIENT_ID
	clientSecret := global.Config.Oauth2Google.CLIENT_SECRET
	domainBe := global.Config.Server.Domain
	clientCallbackURL := fmt.Sprintf("%s/v1/2024/auth/google/callback", domainBe)
	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL, "email", "profile"),
	)
}
