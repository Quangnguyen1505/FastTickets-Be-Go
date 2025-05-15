package services

import (
	"context"

	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
)

type (
	IEvent interface {
		GetAllEventsActive(ctx context.Context, query model.EventQuery) (resultCode int, out []database.GetAllActiveEventsWithLikesRow, err error)
		NewEvent(ctx context.Context, userId string, in *model.AddNewEventParams) (resultCode int, out database.PreGoEvent, err error)
		EditEvent(ctx context.Context, id string, in *model.UpdateEventParams) (resultCode int, err error)
		GetEventById(ctx context.Context, id string) (resultCode int, out database.PreGoEvent, err error)
		DeleteEvent(ctx context.Context, id string) (resultCode int, err error)
		IsLiked(ctx context.Context, userId string) (resultCode int, out []database.GetEventsUserLikeRow, err error)
		EventsLike(ctx context.Context, id string, userId string) (resultCode int, err error)
		EventsUnLike(ctx context.Context, id string, userId string) (resultCode int, err error)
	}
)

var (
	localEvent IEvent
)

func Event() IEvent {
	if localEvent == nil {
		panic("implement localEvent not found for interface IEvent")
	}

	return localEvent
}

func InitEvent(i IEvent) {
	localEvent = i
}
