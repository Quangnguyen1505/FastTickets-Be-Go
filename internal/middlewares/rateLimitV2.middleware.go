package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/utils/context"
	"github.com/ntquang/ecommerce/response"
	limiter "github.com/ulule/limiter/v3"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiterV2 struct {
	globalRateLimiter         *limiter.Limiter
	publicAPIRateLimiter      *limiter.Limiter
	userPrivateAPIRateLimiter *limiter.Limiter
}

func NewRateLimiterV2() *RateLimiterV2 {
	rateLimit := &RateLimiterV2{
		globalRateLimiter:         rateLimiterV2("100-S"),
		publicAPIRateLimiter:      rateLimiterV2("80-S"),
		userPrivateAPIRateLimiter: rateLimiterV2("50-S"),
	}
	return rateLimit
}

func rateLimiterV2(interval string) *limiter.Limiter {
	store, err := redisStore.NewStoreWithOptions(global.Redis, limiter.StoreOptions{
		Prefix:          "rate-limiter", //stand before the key with the purpose of this key that key, ex: rate-limiter:clientip
		MaxRetry:        3,
		CleanUpInterval: time.Hour,
	})
	if err != nil {
		return nil
	}

	rate, err := limiter.NewRateFromFormatted(interval)
	if err != nil {
		panic(err)
	}

	instance := limiter.New(store, rate)
	return instance
}

// GLOBAL RATE LIMITER
func (rl *RateLimiterV2) GlobalRateLimiterV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "global"
		log.Println("global--->")
		limiterContext, err := rl.globalRateLimiter.Get(c, key)
		if err != nil {
			fmt.Println("Failed to check rate limit Global ", err)
			c.Next()
			return
		}
		if limiterContext.Reached {
			log.Printf("Rate limit breached GLOBAL %s", key)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached GLOBAL, try later"})
			return
		}
		c.Next()
	}
}

// PUBLIC API RATE LIMITER
func (rl *RateLimiterV2) PublicAPIRateLimiterV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		rateLimitPath := rl.publicAPIRateLimiter
		log.Println("Client IP ----> ", c.ClientIP()) //ip for client

		key := fmt.Sprintf("%s", c.ClientIP())
		limiterContext, err := rateLimitPath.Get(c, key)
		if err != nil {
			fmt.Println("Failed to check rate limit Public ", err)
			c.Next()
			return
		}
		if limiterContext.Reached {
			log.Printf("Rate limit breached Public %s", key)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached Public, try later"})
			return
		}
		c.Next()
	}
}

// USER PRIVATE API RATE LIMITER
func (rl *RateLimiterV2) UserPrivateAPIRateLimiterV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		rateLimitPath := rl.userPrivateAPIRateLimiter

		userId, err := context.GetUserIdFromUUID(c.Request.Context())
		if err != nil {
			response.ErrorResponse(c, response.ErrTwoFactorAuthSetUpFailed, "Missing get UUID", err)
		}
		fmt.Println("userId ", userId)
		key := fmt.Sprintf("%s-%s", userId, c.ClientIP())
		limiterContext, err := rateLimitPath.Get(c, key)
		if err != nil {
			fmt.Println("Failed to check rate limit User Private ", err)
			c.Next()
			return
		}
		if limiterContext.Reached {
			log.Printf("Rate limit breached User Private %s", key)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached User Private, try later"})
			return
		}
		c.Next()
	}
}
