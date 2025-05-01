package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/ntquang/ecommerce/global"
	consts "github.com/ntquang/ecommerce/internal/const"
	"github.com/ntquang/ecommerce/internal/helper"
	"github.com/ntquang/ecommerce/internal/utils/sendto"
	amqp "github.com/rabbitmq/amqp091-go"
)

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

	StartAllConsumers()
}

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

type ResetPasswordEmailMessage struct {
	Email    string `json:"email"`
	LinkHtml string `json:"linkHtml"`
}

func StartAllConsumers() {
	helper.StartConsumer(
		"email_exchange",
		"direct",
		"email_queue",
		"booking.success",
		func(body []byte) error {
			var emailMsg EmailMessage
			if err := json.Unmarshal(body, &emailMsg); err != nil {
				return err
			}
			return sendto.SendTemplateEmailOtp(
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
		},
	)

	helper.StartConsumer(
		"email_exchange",
		"direct",
		"reset_password_queue",
		"reset.password",
		func(body []byte) error {
			var resetMsg ResetPasswordEmailMessage
			if err := json.Unmarshal(body, &resetMsg); err != nil {
				return err
			}
			fmt.Println("linkHtml", resetMsg.LinkHtml)
			return sendto.SendTemplateEmailOtp(
				[]string{resetMsg.Email},
				consts.EMAIL_SEND,
				"reset-password.html",
				map[string]interface{}{
					"linkHtml": resetMsg.LinkHtml,
				},
				"Đặt lại mật khẩu",
			)
		},
	)
}
