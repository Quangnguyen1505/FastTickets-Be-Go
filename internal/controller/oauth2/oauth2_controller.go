package oauth2

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"net/http"
)

var Oauth2 = new(cUserLoginOauth2)

type cUserLoginOauth2 struct{}

func (oauth2 *cUserLoginOauth2) LoginWithProvider(ctx *gin.Context) {
	provider := ctx.Param("provider")
	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (oauth2 *cUserLoginOauth2) CallbackHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	accessToken := user.AccessToken
	refreshToken := user.RefreshToken
	expiresAt := user.ExpiresAt

	// In thông tin ra
	fmt.Fprintf(ctx.Writer, "Access Token: %s\n", accessToken)
	fmt.Fprintf(ctx.Writer, "Refresh Token: %s\n", refreshToken)
	fmt.Fprintf(ctx.Writer, "Expires At: %s\n", expiresAt.String())
	fmt.Fprintf(ctx.Writer, "User Info: %s - %s\n", user.Name, user.Email)
}
