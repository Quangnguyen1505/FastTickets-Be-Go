package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractBearerToken(c *gin.Context) (string, bool) {
	//Authorization: Bearer token
	authHeader := c.GetHeader("authorization")
	if strings.HasPrefix(authHeader, "Bearer") {
		return strings.TrimPrefix(authHeader, "Bearer "), true
	}

	return "", true
}
