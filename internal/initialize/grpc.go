package initialize

import (
	"log"

	"github.com/ntquang/ecommerce/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGrpcClient() {
	address := "localhost:8083"
	// Kết nối đến server gRPC (không cần chứng chỉ TLS).
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("didn't connect: %v ", err)
	}
	// defer conn.Close()

	global.Grpc = conn
	global.Logger.Info("Grpc initialized successfully!")
}
