package helper

import (
	"fmt"

	"github.com/ntquang/ecommerce/global"
	"go.uber.org/zap"
)

func StartConsumer(exchangeName, exchangeType, queueName, routingKey string, handlerFunc func([]byte) error) {
	err := global.RabbitMQChannel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		global.Logger.Error("Failed to declare exchange:", zap.Error(err))
		return
	}

	_, err = global.RabbitMQChannel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		global.Logger.Error("Failed to declare queue:", zap.Error(err))
		return
	}

	err = global.RabbitMQChannel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("Failed to bind queue:", zap.Error(err))
		return
	}

	msgs, err := global.RabbitMQChannel.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		global.Logger.Error("Failed to start consumer", zap.Error(err))
		return
	}

	go func() {
		for msg := range msgs {
			if err := handlerFunc(msg.Body); err != nil {
				global.Logger.Error("Handler error:", zap.Error(err))
			}
		}
	}()

	global.Logger.Info(fmt.Sprintf("Consumer started for queue: %s", queueName))
}
