package imple

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/pkg/errors"
)

type ChatMessage struct {
	SenderID string `json:"sender_id"`
	Message  string `json:"message"`
	SentAt   string `json:"sent_at"` // ISO8601 format
}

func InitSession(ctx context.Context, userId string) (resultCode int, out string /* sessionId */, err error) {
	sessionId := uuid.NewString()

	sessionKey := fmt.Sprintf("chat:session:%s", sessionId)
	userSessionsKey := fmt.Sprintf("chat:sessions:user:%s", userId)

	now := time.Now().Format(time.RFC3339)

	sessionData := map[string]interface{}{
		"session_id":  sessionId,
		"customer_id": userId,
		"staff_id":    "", // chưa có staff join
		"created_at":  now,
		"updated_at":  now,
		"status":      "open",
	}

	err = global.Redis.HSet(ctx, sessionKey, sessionData).Err()
	if err != nil {
		resultCode = 500
		return
	}

	err = global.Redis.LPush(ctx, userSessionsKey, sessionId).Err()
	if err != nil {
		resultCode = 500
		return
	}

	resultCode = 200
	out = sessionId
	return
}

func GetHistoryChat(ctx context.Context, sessionId string) (resultCode int, out []model.GetHistoryChatRow, err error) {
	key := fmt.Sprintf("chat:messages:%s", sessionId)

	// Lấy toàn bộ tin nhắn
	messages, err := global.Redis.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		resultCode = 500
		return
	}

	for _, msgStr := range messages {
		var msg model.GetHistoryChatRow
		if err = json.Unmarshal([]byte(msgStr), &msg); err != nil {
			resultCode = 500
			return
		}
		out = append(out, msg)
	}

	resultCode = 200
	return
}

func SendMessage(ctx context.Context, sessionId string, sender string, message string) (resultCode int, err error) {
	key := fmt.Sprintf("chat:messages:%s", sessionId)

	msg := ChatMessage{
		SenderID: sender,
		Message:  message,
		SentAt:   time.Now().Format(time.RFC3339),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		resultCode = 500
		return
	}

	// Đẩy tin nhắn vào Redis List
	err = global.Redis.RPush(ctx, key, msgBytes).Err()
	if err != nil {
		resultCode = 500
		return
	}

	sessionKey := fmt.Sprintf("chat:session:%s", sessionId)
	global.Redis.HMSet(ctx, sessionKey, map[string]interface{}{
		"last_message": message,
		"updated_at":   time.Now().Format(time.RFC3339),
	})

	// send real time here

	// Optional: Giới hạn độ dài list để tránh quá tải
	global.Redis.LTrim(ctx, key, -100, -1) // Giữ lại 100 tin mới nhất

	resultCode = 200
	return
}

func GetAllSession(ctx context.Context, userId string) (resultCode int, out []model.GetAllSessionRow, err error) {
	userSessionsKey := fmt.Sprintf("chat:user:%s:sessions", userId)

	// Lấy tất cả sessionId của user (có thể là customer hoặc staff)
	sessionIds, err := global.Redis.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		resultCode = 500
		err = errors.Wrap(err, "redis SMembers error")
		return
	}

	out = make([]model.GetAllSessionRow, 0, len(sessionIds))

	for _, sessionId := range sessionIds {
		sessionKey := fmt.Sprintf("chat:session:%s", sessionId)

		// Lấy dữ liệu session dạng map
		sessionData, err := global.Redis.HGetAll(ctx, sessionKey).Result()
		if err != nil {
			continue // lỗi Redis thì bỏ qua session này
		}
		if len(sessionData) == 0 {
			continue // session không tồn tại
		}

		// Parse các trường
		customerID := sessionData["customer_id"]
		staffIDStr := sessionData["staff_id"]
		var staffID *string
		if staffIDStr != "" {
			staffID = &staffIDStr
		}

		createdAtStr := sessionData["created_at"]
		createdAt, _ := time.Parse(time.RFC3339, createdAtStr)

		lastMessage := sessionData["last_message"]
		if lastMessage == "" {
			lastMessage = "" // hoặc nil
		}
		var lastMsgPtr *string
		if lastMessage != "" {
			lastMsgPtr = &lastMessage
		}

		updatedAtStr := sessionData["updated_at"]
		updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

		status := sessionData["status"]
		if status == "" {
			status = "open" // mặc định
		}

		out = append(out, model.GetAllSessionRow{
			SessionID:   sessionId,
			CustomerID:  customerID,
			StaffID:     staffID,
			CreatedAt:   createdAt,
			LastMessage: lastMsgPtr,
			UpdatedAt:   updatedAt,
			Status:      status,
		})
	}

	resultCode = 200
	return
}

func JoinChatSession(ctx context.Context, sessionId string, staffId string) (resultCode int, err error) {
	sessionKey := fmt.Sprintf("chat:session:%s", sessionId)
	staffSessionsKey := fmt.Sprintf("chat:user:%s:sessions", staffId)

	// Kiểm tra session có tồn tại không
	exists, err := global.Redis.Exists(ctx, sessionKey).Result()
	if err != nil {
		resultCode = 500
		return
	}
	if exists == 0 {
		resultCode = 404
		err = fmt.Errorf("session not found")
		return
	}

	// Gán staffId vào session
	if err = global.Redis.HSet(ctx, sessionKey, "staff_id", staffId).Err(); err != nil {
		resultCode = 500
		return
	}

	// Thay đổi trạng thái session thành "in_progress"
	if err = global.Redis.HSet(ctx, sessionKey, "status", "in_progress").Err(); err != nil {
		resultCode = 500
		return
	}

	// Thêm sessionId vào danh sách sessions của staff
	if err = global.Redis.SAdd(ctx, staffSessionsKey, sessionId).Err(); err != nil {
		resultCode = 500
		return
	}

	resultCode = 200
	return
}

func CloseChatSession(ctx context.Context, sessionId string) (resultCode int, err error) {
	sessionKey := fmt.Sprintf("chat:session:%s", sessionId)

	// Kiểm tra session có tồn tại không
	exists, err := global.Redis.Exists(ctx, sessionKey).Result()
	if err != nil {
		resultCode = 500
		return
	}
	if exists == 0 {
		resultCode = 404
		err = fmt.Errorf("session not found")
		return
	}

	// Thay đổi trạng thái session thành "closed"
	if err = global.Redis.HSet(ctx, sessionKey, "status", "closed").Err(); err != nil {
		resultCode = 500
		return
	}

	resultCode = 200
	return
}
