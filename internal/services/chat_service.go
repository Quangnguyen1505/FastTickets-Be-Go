package services

import (
	"context"

	"github.com/ntquang/ecommerce/internal/model"
)

type (
	IChat interface {
		// 1. Tạo hoặc lấy lại session cho khách hàng (ẩn danh hoặc đã login)
		InitSession(ctx context.Context, userId string) (resultCode int, out string /* sessionId */, err error)

		// 2. Lấy lịch sử tin nhắn theo sessionId (từ Redis hoặc DB nếu bạn có backup lâu dài)
		GetHistoryChat(ctx context.Context, sessionId string) (resultCode int, out []model.GetHistoryChatRow, err error)

		// 3. Gửi tin nhắn vào session
		// Gợi ý: Nên thêm trường sender ("customer"/"staff") để phân biệt phía gửi
		SendMessage(ctx context.Context, in *model.ChatMessageParams) (resultCode int, err error)

		// 4. Lấy danh sách tất cả session mà nhân viên có thể thấy (gợi ý: có thể filter theo trạng thái open/closed)
		GetAllSession(ctx context.Context) (resultCode int, out []model.GetAllSessionRow, err error)

		// 5. Nhân viên tham gia vào một phiên chat
		JoinChatSession(ctx context.Context, in *model.JoinChatSessionParams) (resultCode int, err error)

		// 6. (Tuỳ chọn) Đóng phiên chat
		CloseChatSession(ctx context.Context, sessionId string) (resultCode int, err error)
	}
)

var (
	localChat IChat
)

func Chat() IChat {
	if localChat == nil {
		panic("implement localChat not found for interface IChat")
	}

	return localChat
}

func InitChat(i IChat) {
	localChat = i
}
