package initialize

import (
	"fmt"

	"github.com/ntquang/ecommerce/global"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(fmt.Errorf("failed to start conn rabbitmq: %w", err))
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		panic(fmt.Errorf("failed to create channel rabbitmq: %w", err))
	}

	global.RabbitMQChannel = channel
	global.Logger.Info("RabbitMQ initialized successfully!")
}
