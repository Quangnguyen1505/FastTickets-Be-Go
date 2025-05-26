package imple

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/internal/websocket"
	"github.com/pkg/errors"
)

type sChatEmployee struct{}

func NewChatEmployee() *sChatEmployee {
	return &sChatEmployee{}
}

func (s *sChatEmployee) InitSession(ctx context.Context, userId string) (resultCode int, out string /* sessionId */, err error) {
	sessionId := uuid.NewString()

	sessionKey := fmt.Sprintf("chat:session:%s", sessionId)
	userSessionsKey := fmt.Sprintf("chat:sessions:user:%s", userId)

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now().In(loc).Format(time.RFC3339)

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

	global.Redis.SAdd(ctx, "chat:sessions:all", sessionId)

	err = global.Redis.LPush(ctx, userSessionsKey, sessionId).Err()
	if err != nil {
		resultCode = 500
		return
	}

	global.Redis.LTrim(ctx, userSessionsKey, 0, 99)

	resultCode = 200
	out = sessionId
	return
}

func (s *sChatEmployee) GetHistoryChat(ctx context.Context, sessionId string) (resultCode int, out []model.GetHistoryChatRow, err error) {
	key := fmt.Sprintf("chat:messages:%s", sessionId)

	// Lấy toàn bộ tin nhắn
	messages, err := global.Redis.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		resultCode = 500
		return
	}

	println("messages", messages)

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

func (s *sChatEmployee) SendMessage(ctx context.Context, in *model.ChatMessageParams) (resultCode int, err error) {
	key := fmt.Sprintf("chat:messages:%s", in.SessionId)

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now().In(loc).Format(time.RFC3339)

	msg := model.ChatMessage{
		SenderID: in.Sender,
		Message:  in.Message,
		SentAt:   now,
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

	sessionKey := fmt.Sprintf("chat:session:%s", in.SessionId)
	global.Redis.HMSet(ctx, sessionKey, map[string]interface{}{
		"last_message": in.Message,
		"updated_at":   now,
	})

	// send real time here
	websocket.ChatHub.Broadcast <- websocket.BroadcastMsg{
		SessionID: in.SessionId,
		Message:   msgBytes,
	}

	// Optional: Giới hạn độ dài list để tránh quá tải
	global.Redis.LTrim(ctx, key, -100, -1) // Giữ lại 100 tin mới nhất

	resultCode = 200
	return
}

func (s *sChatEmployee) GetAllSession(ctx context.Context) (resultCode int, out []model.GetAllSessionRow, err error) {
	allSessionsKey := "chat:sessions:all" // Danh sách tất cả session ID trong hệ thống

	sessionIds, err := global.Redis.SMembers(ctx, allSessionsKey).Result()
	if err != nil {
		resultCode = 500
		err = errors.Wrap(err, "redis SMembers error")
		return
	}

	out = make([]model.GetAllSessionRow, 0)

	for _, sessionId := range sessionIds {
		sessionKey := fmt.Sprintf("chat:session:%s", sessionId)

		sessionData, err := global.Redis.HGetAll(ctx, sessionKey).Result()
		if err != nil || len(sessionData) == 0 {
			continue
		}

		// Chỉ lấy các session chưa có staff xử lý
		// if sessionData["staff_id"] != "" || sessionData["status"] != "open" {
		// 	continue
		// }

		status := sessionData["status"]
		if status != "open" && status != "in_progress" {
			continue
		}

		customerID := sessionData["customer_id"]

		var staffID *string = nil

		createdAt, _ := time.Parse(time.RFC3339, sessionData["created_at"])
		updatedAt, _ := time.Parse(time.RFC3339, sessionData["updated_at"])

		var lastMsgPtr *string
		lastMessage := sessionData["last_message"]
		if lastMessage != "" {
			lastMsgPtr = &lastMessage
		}

		out = append(out, model.GetAllSessionRow{
			SessionID:   sessionId,
			CustomerID:  customerID,
			StaffID:     staffID,
			CreatedAt:   createdAt,
			LastMessage: lastMsgPtr,
			UpdatedAt:   updatedAt,
			Status:      sessionData["status"],
		})
	}

	resultCode = 200
	return
}

// dùng sessionID của khách hàng chứ không phải của admin
// vậy admin không cần init session
func (s *sChatEmployee) JoinChatSession(ctx context.Context, in *model.JoinChatSessionParams) (resultCode int, err error) {
	sessionKey := fmt.Sprintf("chat:session:%s", in.SessionId)
	staffSessionsKey := fmt.Sprintf("chat:user:%s:sessions", in.StaffId)

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
	if err = global.Redis.HSet(ctx, sessionKey, "staff_id", in.StaffId).Err(); err != nil {
		resultCode = 500
		return
	}

	// Thay đổi trạng thái session thành "in_progress"
	if err = global.Redis.HSet(ctx, sessionKey, "status", "in_progress").Err(); err != nil {
		resultCode = 500
		return
	}

	// Thêm sessionId vào danh sách sessions của staff
	if err = global.Redis.SAdd(ctx, staffSessionsKey, in.SessionId).Err(); err != nil {
		resultCode = 500
		return
	}

	joinEvent := map[string]interface{}{
		"event":      "staff_joined",
		"session_id": in.SessionId,
		"staff_id":   in.StaffId,
	}
	msgBytes, _ := json.Marshal(joinEvent)
	websocket.ChatHub.Broadcast <- websocket.BroadcastMsg{
		SessionID: in.SessionId,
		Message:   msgBytes,
	}

	resultCode = 200
	return
}

func (s *sChatEmployee) CloseChatSession(ctx context.Context, sessionId string) (resultCode int, err error) {
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

	message := map[string]interface{}{
		"event":      "session_closed",
		"session_id": sessionId,
		"message":    "Vấn đề của bạn đã được giải quyết. Bạn có còn điều gì muốn hỏi không?",
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		resultCode = 500
		return
	}

	websocket.ChatHub.Broadcast <- websocket.BroadcastMsg{
		SessionID: sessionId,
		Message:   msgBytes,
	}

	resultCode = 200
	return
}
