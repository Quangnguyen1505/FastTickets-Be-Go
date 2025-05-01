package middlewares

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	consts "github.com/ntquang/ecommerce/internal/const"
	"github.com/ntquang/ecommerce/internal/utils/auth"
	"github.com/ntquang/ecommerce/internal/utils/cache"
	"github.com/ntquang/ecommerce/response"
)

type UserTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	PublicKey    string `json:"publicKey"`
	RoleId       string `json:"roleId"`
}

func Authentication() gin.HandlerFunc {
	// Add authentication logic here
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		log.Println("Url request ::", url)

		//get token in header
		JwtToken, ok := auth.ExtractBearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"code":        40001,
				"err":         "Unauthorized",
				"description": "",
			})
			return
		}
		fmt.Println("accesstoken header ", JwtToken)
		userId := c.GetHeader("x-client-id")
		keyUser := fmt.Sprintf("%s%s", consts.CACHE_USER, userId)

		var userToken UserTokens
		if err := cache.GetHashCache(c.Request.Context(), keyUser, &userToken); err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"code":        response.ErrCodeRedisGetFailed,
				"err":         "get hash cache",
				"description": "",
			})
			return
		}

		accessToken := userToken.AccessToken
		fmt.Println("accesstoken redis ", accessToken)
		if accessToken != JwtToken {
			c.AbortWithStatusJSON(500, gin.H{
				"code":        response.ErrCodeParamInvalid,
				"err":         "accessToken",
				"description": "accessToken not match with db",
			})
			return
		}

		//validate jwt token by subject
		publicKey := userToken.PublicKey
		claims, err := auth.VerifyToken(JwtToken, publicKey)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":        40001,
				"err":         "Invalid Token",
				"description": "" + err.Error(),
			})
			return
		}

		log.Println("Claims UUID:: ", claims.Subject) //userid
		// Tạo context mới và thêm cả subjectUUID và roleId vào đó
		ctx := context.WithValue(c.Request.Context(), "subjectUUID", claims.Subject)
		ctx = context.WithValue(ctx, "roleId", fmt.Sprintf("%v", claims.RoleId))

		// Cập nhật context trong request
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func CheckPermission() gin.HandlerFunc {
	// Add authentication logic here
	return func(c *gin.Context) {
		roleId, ok := c.Request.Context().Value("roleId").(string)
		fmt.Println("roleId", roleId)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"code":        40001,
				"err":         "Missing Role",
				"description": "Role ID not found in context",
			})
			return
		}

		userId := c.GetHeader("x-client-id")
		keyUser := fmt.Sprintf("%s%s", consts.CACHE_USER, userId)

		value, err := cache.GetFieldHashCache(c.Request.Context(), keyUser, "roleId")
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"code":        response.ErrCodeRedisGetFailed,
				"err":         "get hash cache",
				"description": "",
			})
			return
		}
		fmt.Println("roleId redis ", value)
		fmt.Println("roleId context ", roleId)
		if value != roleId {
			c.AbortWithStatusJSON(401, gin.H{
				"code":        40001,
				"err":         "Unauthorized",
				"description": "Permission denied",
			})
			return
		}
		c.Next()
	}
}
