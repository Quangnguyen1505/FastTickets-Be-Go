package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	consts "github.com/ntquang/ecommerce/internal/const"
)

func ExtractBearerToken(c *gin.Context) (string, bool) {
	//Authorization: Bearer token
	authHeader := c.GetHeader(consts.Authorization)
	if strings.HasPrefix(authHeader, "Bearer") {
		return strings.TrimPrefix(authHeader, "Bearer "), true
	}

	return "", true
}
