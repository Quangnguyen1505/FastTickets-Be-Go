package chatemployee

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/internal/services"
	"github.com/ntquang/ecommerce/internal/utils/context"
	"github.com/ntquang/ecommerce/response"
	"go.uber.org/zap"
)

var ChatEmployee = new(CChatEmployee)

type CChatEmployee struct{}

func (cCEmpl *CChatEmployee) InitSessionId(ctx *gin.Context) {
	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrTwoFactorAuthSetUpFailed, "Missing get UUID", err)
	}

	statusCode, sessionId, err := services.Chat().InitSession(ctx, userId)
	if err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params init sessionID", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, sessionId, "New Session Created Successfully!")
}

func (cCEmpl *CChatEmployee) GetHistoryChat(ctx *gin.Context) {
	sessionId := ctx.Param("sessionId")
	if sessionId == "" {
		global.Logger.Error("Error parse params get history chat", zap.Error(fmt.Errorf("sessionId is required")))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf("sessionId is required"))
		return
	}

	statusCode, history, err := services.Chat().GetHistoryChat(ctx, sessionId)
	if err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params get history chat", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, history, "Get History Chat Successfully!")
}

func (cCEmpl *CChatEmployee) SendMessage(ctx *gin.Context) {
	var params model.ChatMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params send message", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	statusCode, err := services.Chat().SendMessage(ctx, &params)
	if err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params send message", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, nil, "Send Message Successfully!")
}

func (cCEmpl *CChatEmployee) GetAllSessionId(ctx *gin.Context) {
	statusCode, sessionIds, err := services.Chat().GetAllSession(ctx)
	if err != nil {
		fmt.Println("Lỗi:", err)
		global.Logger.Error("Error parse params get all session ID", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, sessionIds, "Get All Session ID Successfully!")
}

func (cCEmpl *CChatEmployee) JoinChatSession(ctx *gin.Context) {
	var params model.JoinChatSessionParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params join chat session", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	statusCode, err := services.Chat().JoinChatSession(ctx, &params)
	if err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params join chat session", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, nil, "Join Chat Session Successfully!")
}

func (cCEmpl *CChatEmployee) CloseChatSession(ctx *gin.Context) {
	sessionId := ctx.Param("sessionId")
	if sessionId == "" {
		global.Logger.Error("Error parse params close chat session", zap.Error(fmt.Errorf("sessionId is required")))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf("sessionId is required"))
		return
	}

	statusCode, err := services.Chat().CloseChatSession(ctx, sessionId)
	if err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params close chat session", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, nil, "Close Chat Session Successfully!")
}
