package services

import (
	"context"

	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
)

type (
	IContactMessage interface {
		GetAllContactMessageByStatus(ctx context.Context, query model.ContactMessageQuery) (resultCode int, out []database.PreGoContactMessage, err error)
		NewContactMessage(ctx context.Context, in *model.AddNewContactMessageParams) (resultCode int, out database.PreGoContactMessage, err error)
		EditStatusContactMessage(ctx context.Context, id string, status int16) (resultCode int, err error)
		GetContactMessageById(ctx context.Context, id string) (resultCode int, out database.PreGoContactMessage, err error)
		DeleteContactMessage(ctx context.Context, id string) (resultCode int, err error)
		SendEmailToCustomer(ctx context.Context, in *model.ResponseCustomer) (resultCode int, err error)
	}
)

var (
	localContactMessage IContactMessage
)

func ContactMessage() IContactMessage {
	if localContactMessage == nil {
		panic("implement localContactMessage not found for interface IContactMessage")
	}

	return localContactMessage
}

func InitContactMessage(i IContactMessage) {
	localContactMessage = i
}
