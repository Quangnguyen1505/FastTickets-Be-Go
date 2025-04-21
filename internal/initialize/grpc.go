package initialize

import (
	"fmt"
	"log"

	"github.com/ntquang/ecommerce/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGrpcClient() {
	address := fmt.Sprintf("%s:%d", global.Config.Grpc.Client.Host, global.Config.Grpc.Client.Port)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("didn't connect: %v ", err)
	}
	// defer conn.Close()

	global.Grpc = conn
	global.Logger.Info("Grpc initialized successfully!")
}
