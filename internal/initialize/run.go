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
	StartConsumerExDirect()
	InitOauth2()
	r := Initrouter()

	return r
}
