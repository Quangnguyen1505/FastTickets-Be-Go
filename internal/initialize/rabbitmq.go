package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/ntquang/ecommerce/global"
	consts "github.com/ntquang/ecommerce/internal/const"
	"github.com/ntquang/ecommerce/internal/utils/sendto"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type SeatDetail struct {
	Seat     string `json:"seat"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	STT      int    `json:"stt"`
}

type EmailMessage struct {
	Email      string       `json:"email"`
	BookingID  string       `json:"booking_id"`
	MovieTitle string       `json:"movie_title"`
	AgeRating  string       `json:"age_rating"`
	RoomName   string       `json:"room_name"`
	ShowDate   string       `json:"show_date"`
	StartTime  string       `json:"start_time"`
	SeatNumber int          `json:"seat_number"`
	Booking    []SeatDetail `json:"booking"`
	TotalPrice int          `json:"total_price"`
}

func InitRabbitMQ() {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", global.Config.RabbitMQ.Username, global.Config.RabbitMQ.Password, global.Config.RabbitMQ.Host, global.Config.RabbitMQ.Port)
	conn, err := amqp.Dial(uri)
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

func StartConsumerExDirect() {
	queueName := "email_queue"
	exchange := "email_exchange"
	routingKey := "booking.success"

	err := global.RabbitMQChannel.ExchangeDeclare(
		exchange, // Tên exchange
		"direct", // Kiểu exchange: direct, fanout, topic, headers
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		global.Logger.Error("Failed to declare exchange:", zap.Error(err))
		return
	}

	_, err = global.RabbitMQChannel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		global.Logger.Error("Failed to declare queue:", zap.Error(err))
		return
	}

	err = global.RabbitMQChannel.QueueBind(
		queueName,  // Queue name
		routingKey, // Routing key
		exchange,   // Exchange name
		false,      // No-wait
		nil,        // Arguments
	)
	if err != nil {
		global.Logger.Error("Failed to bind queue:", zap.Error(err))
		return
	}

	msgs, err := global.RabbitMQChannel.Consume(
		queueName, // Queue name
		"",        // Consumer tag
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		global.Logger.Error("Failed to start consumer", zap.Error(err))
		return
	}

	// Process messages trong một goroutine
	go func() {
		for msg := range msgs {
			// Xử lý message tại đây (ví dụ gửi email)
			// global.Logger.Info("Received message", zap.String("message", string(msg.Body)))
			// nhận message dưới dạng json và parse json
			var emailMsg EmailMessage
			if err := json.Unmarshal(msg.Body, &emailMsg); err != nil {
				global.Logger.Error("Failed to parse message", zap.Error(err))
				continue
			}

			err = sendto.SendTemplateEmailOtp(
				[]string{emailMsg.Email},
				consts.EMAIL_SEND,
				"booking-success.html",
				map[string]interface{}{
					"ticket_code": emailMsg.BookingID,
					"movie_title": emailMsg.MovieTitle,
					"age_rating":  emailMsg.AgeRating,
					"show_date":   emailMsg.ShowDate,
					"start_time":  emailMsg.StartTime,
					"room_name":   emailMsg.RoomName,
					"total_seats": emailMsg.SeatNumber,
					"SeatDetails": emailMsg.Booking,
					"TotalPrice":  emailMsg.TotalPrice,
				},
				"Thông tin đặt vé xem phim",
			)
			// Nếu bạn không sử dụng auto-ack, hãy nhớ call Ack sau khi xử lý xong message
			// msg.Ack(false)
		}
	}()

	global.Logger.Info("Consumer started and listening to messages...")
}
