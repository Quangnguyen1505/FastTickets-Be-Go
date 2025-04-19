package routers

// import (
// 	"github.com/gin-gonic/gin"
// 	c "github.com/ntquang/ecommerce/internal/controller"
// 	"github.com/ntquang/ecommerce/internal/middlewares"
// )

// func NewRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.Use(middlewares.Authentication())
// 	v1 := r.Group("/v1")
// 	{
// 		v1.GET("/ping", c.NewPongController().Pong)
// 		v1.GET("/user/1", c.NewUserController().GetUser)
// 		// v1.POST("/ping", Pong)
// 		// v1.PATCH("/ping", Pong)
// 		// v1.DELETE("/ping", Pong)
// 		// v1.PUT("/ping", Pong)
// 		// v1.OPTIONS("/ping", Pong)
// 	}

// 	// v2 := r.Group("/v2")
// 	// {
// 	// 	v2.GET("/ping", Pong)
// 	// 	v2.POST("/ping", Pong)
// 	// 	v2.PATCH("/ping", Pong)
// 	// 	v2.DELETE("/ping", Pong)
// 	// 	v2.PUT("/ping", Pong)
// 	// 	v2.OPTIONS("/ping", Pong)
// 	// }

// 	return r
// }
