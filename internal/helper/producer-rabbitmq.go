package helper

import (
	"encoding/json"

	"github.com/ntquang/ecommerce/global"
	"github.com/rabbitmq/amqp091-go"
)

func SendToRabbitMQ(message interface{}, exchange, routingKey string) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = global.RabbitMQChannel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
