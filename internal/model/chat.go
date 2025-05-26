package model

import "time"

type GetHistoryChatRow struct {
	Sender    string `json:"sender_id"` // "customer" hoặc "staff"
	Message   string `json:"message"`
	Timestamp string `json:"sent_at"` // Có thể lưu theo time.Time hoặc int64 Unix timestamp
}

type GetAllSessionRow struct {
	SessionID   string    `json:"session_id"`
	CustomerID  string    `json:"customer_id"`
	StaffID     *string   `json:"staff_id"` // Có thể là nil nếu chưa ai join
	CreatedAt   time.Time `json:"created_at"`
	LastMessage *string   `json:"last_message"` // Có thể là preview tin nhắn cuối
	UpdatedAt   time.Time `json:"updated_at"`   // Thời gian tin nhắn cuối
	Status      string    `json:"status"`       // "open", "in_progress", "closed"
}

type ChatMessageParams struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	SessionId string `json:"sessionId"`
}

type ChatMessage struct {
	SenderID string `json:"sender_id"`
	Message  string `json:"message"`
	SentAt   string `json:"sent_at"` // ISO8601 format
}

type JoinChatSessionParams struct {
	SessionId string `json:"sessionId"`
	StaffId   string `json:"staffId"`
}
