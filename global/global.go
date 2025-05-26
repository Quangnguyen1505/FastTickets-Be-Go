package global

import (
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntquang/ecommerce/pkg/logger"
	"github.com/ntquang/ecommerce/pkg/setting"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var (
	Config          setting.Config
	Logger          *logger.LoggerZap
	Pdbc            *pgxpool.Pool
	Redis           *redis.Client
	RabbitMQChannel *amqp.Channel
	Grpc            *grpc.ClientConn
	SocketHub       *websocket.Conn
)
