package main

import (
	"fmt"

	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/initialize"
	ginSwagger "github.com/swaggo/gin-swagger"

	// gin-swagger middlewares
	_ "github.com/ntquang/ecommerce/docs"

	swaggerFiles "github.com/swaggo/files"
)

// @title           Demo API Ecommerce
// @version         1.0.0
// @description     This is a server ecommerce.
// @termsOfService  https://github.com/Quangnguyen1505/Ecommerce-Go

// @contact.name   TEAM QUANG
// @contact.url    https://github.com/Quangnguyen1505/Ecommerce-Go
// @contact.email  quang0706r@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /v1/2024
func main() {
	r := initialize.Run()
	port := global.Config.Server.Port

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf(":%v", port))
	global.Logger.Info(fmt.Sprintf("Server running is port %d", port))
}
