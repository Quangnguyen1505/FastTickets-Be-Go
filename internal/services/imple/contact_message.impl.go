package imple

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
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
func (s *sContactMessage) GetAllContactMessageByStatus(ctx context.Context, query model.ContactMessageQuery) (resultCode int, out []database.ContactMessage, err error) {
	limit, page, status := query.Limit, query.Page, query.Status
	offset := (page - 1) * limit
	contactMessages, err := s.r.GetAllContactMessages(ctx, database.GetAllContactMessagesParams{
		Status: int16(status),
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		return response.ErrorListFailed, nil, err
	}

	if len(contactMessages) == 0 {
		return response.ErrorDataNotExists, nil, nil
	}

	return 200, contactMessages, nil
}

func (s *sContactMessage) NewContactMessage(ctx context.Context, in *model.AddNewContactMessageParams) (resultCode int, out database.ContactMessage, err error) {
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

func (s *sContactMessage) EditStatusContactMessage(ctx context.Context, id string, in *model.UpdateContactMessageParams) (resultCode int, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, err
	}
	err = s.r.UpdateContactMessageStatus(ctx, database.UpdateContactMessageStatusParams{
		ID:     uuidID,
		Status: int16(in.Status),
	})

	if err != nil {
		return response.ErrorUpdate, err
	}

	return 200, nil
}

func (s *sContactMessage) GetContactMessageById(ctx context.Context, id string) (resultCode int, out database.ContactMessage, err error) {
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
