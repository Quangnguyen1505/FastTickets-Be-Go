package imple

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	consts "github.com/ntquang/ecommerce/internal/const"
	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/internal/utils/sendto"
	"github.com/ntquang/ecommerce/response"
)

type sContactMessage struct {
	r *database.Queries
}

func NewContactMessage(r *database.Queries) *sContactMessage {
	return &sContactMessage{
		r: r,
	}
}

// implement
func (s *sContactMessage) GetAllContactMessageByStatus(ctx context.Context, query model.ContactMessageQuery) (resultCode int, out []database.PreGoContactMessage, err error) {
	limit, page := query.Limit, query.Page
	offset := (page - 1) * limit

	status := int16(-1)
	if query.Status != nil {
		status = *query.Status
	}

	fmt.Println("status", status)

	contactMessages, err := s.r.GetAllContactMessages(ctx, database.GetAllContactMessagesParams{
		Column1: status,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		return response.ErrorListFailed, nil, err
	}

	if len(contactMessages) == 0 {
		return response.ErrorDataNotExists, nil, nil
	}

	return 200, contactMessages, nil
}

func (s *sContactMessage) NewContactMessage(ctx context.Context, in *model.AddNewContactMessageParams) (resultCode int, out database.PreGoContactMessage, err error) {
	contactMessage, err := s.r.CreateContactMessage(ctx, database.CreateContactMessageParams{
		Name:    in.Name,
		Email:   in.Email,
		Message: in.Message,
		Phone:   pgtype.Text{String: in.Phone, Valid: true},
		Status:  0,
	})

	if err != nil {
		return response.ErrorListFailed, out, err
	}

	return 200, contactMessage, nil
}

func (s *sContactMessage) EditStatusContactMessage(ctx context.Context, id string, status int16) (resultCode int, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, err
	}
	err = s.r.UpdateContactMessageStatus(ctx, database.UpdateContactMessageStatusParams{
		ID:     uuidID,
		Status: int16(status),
	})

	if err != nil {
		return response.ErrorUpdate, err
	}

	return 200, nil
}

func (s *sContactMessage) GetContactMessageById(ctx context.Context, id string) (resultCode int, out database.PreGoContactMessage, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, out, err
	}

	contactMessage, err := s.r.GetContactMessageByID(ctx, uuidID)
	if err != nil {
		return response.ErrorListFailed, out, err
	}

	if len(contactMessage.ID.Bytes) == 0 {
		return response.ErrorDataNotExists, out, nil
	}

	return 200, contactMessage, nil
}

func (s *sContactMessage) DeleteContactMessage(ctx context.Context, id string) (resultCode int, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, err
	}
	contactMessage, err := s.r.GetContactMessageByID(ctx, uuidID)
	if err != nil {
		return response.ErrorDataNotExists, err
	}

	if len(contactMessage.ID.Bytes) == 0 {
		return response.ErrorDataNotExists, nil
	}

	err = s.r.DeleteContactMessage(ctx, uuidID)
	if err != nil {
		return response.ErrorDelete, err
	}

	return 200, nil
}

func (s *sContactMessage) SendEmailToCustomer(ctx context.Context, in *model.ResponseCustomer) (resultCode int, err error) {
	err = sendto.SendTemplateEmailOtp(
		[]string{in.Email},
		consts.EMAIL_SEND,
		"contact-mess.html",
		map[string]interface{}{
			"name":     in.Name,
			"message":  in.Message,
			"response": in.Response,
		},
		"[Phản hồi] Fast Tickets cinema",
	)

	if err != nil {
		return response.ErrorSendEmail, err

	}

	_, err = s.EditStatusContactMessage(ctx, in.ContactId, 2)
	if err != nil {
		return response.ErrorUpdate, err
	}

	return 200, nil
}
