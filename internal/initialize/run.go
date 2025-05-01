package initialize

import "github.com/gin-gonic/gin"

func Run() *gin.Engine {
	InitConfig()
	InitLogger()
	InitPostgresqlC()
	InitRedis()
	InitServiceInterface()
	InitGrpcClient()
	InitRabbitMQ()
	InitOauth2()
	r := Initrouter()

	return r
}
