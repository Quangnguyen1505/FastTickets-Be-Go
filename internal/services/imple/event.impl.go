package imple

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntquang/ecommerce/global"
	consts "github.com/ntquang/ecommerce/internal/const"
	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/helper"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/response"
	"go.uber.org/zap"
)

type sEvent struct {
	r  *database.Queries
	db *pgxpool.Pool
}

func NewEventImpl(r *database.Queries, db *pgxpool.Pool) *sEvent {
	return &sEvent{
		r:  r,
		db: db,
	}
}

// implement
func (evt *sEvent) GetAllEventsActive(ctx context.Context, query model.EventQuery) (resultCode int, out []database.GetAllActiveEventsWithLikesRow, err error) {
	limit, page := query.Limit, query.Page
	offset := (page - 1) * limit

	params := database.GetAllActiveEventsWithLikesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	foundEvent, err := evt.r.GetAllActiveEventsWithLikes(ctx, params)
	if err != nil {
		return response.ErrorListFailed, nil, err
	}

	return 200, foundEvent, nil
}

func (evt *sEvent) NewEvent(ctx context.Context, userId string, in *model.AddNewEventParams) (resultCode int, out database.PreGoEvent, err error) {
	fmt.Println("evenet name and des ", in.Name, in.Description)
	tx, err := evt.db.Begin(ctx)
	if err != nil {
		return response.ErrorTransactionBegin, out, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			commitErr := tx.Commit(ctx)
			if commitErr != nil {
				err = commitErr
			}
		}
	}()

	txQueries := evt.r.WithTx(tx)

	event, err := txQueries.GetEventByName(ctx, in.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return response.ErrorListFailed, out, err
	}
	if event.ID.Valid {
		return response.ErrorListFailed, out, err
	}
	fmt.Println("evenet name and des ", in.Name, in.Description)
	newEvent, err := txQueries.AddNewEvent(ctx, database.AddNewEventParams{
		EventName:        in.Name,
		EventDescription: pgtype.Text{String: in.Description, Valid: true},
		EventImageUrl:    pgtype.Text{String: in.ImageUrl, Valid: true},
		EventStart:       pgtype.Timestamp{Time: in.Start, Valid: true},
		EventEnd:         pgtype.Timestamp{Time: in.End, Valid: true},
		EventActive:      pgtype.Bool{Bool: true, Valid: true},
		UserID:           pgtype.UUID{Bytes: uuid.MustParse(userId), Valid: true},
	})
	if err != nil {
		return response.ErrorInsert, out, nil
	}

	messageNoti := map[string]interface{}{
		"pattern": "noti_created",
		"data": map[string]interface{}{
			"noti_type":    "EVENT",
			"noti_content": "Rạp vừa ra mắt sự kiện mới, xem ngay !",
			"noti_options": map[string]interface{}{
				"id":    newEvent.ID,
				"title": newEvent.EventName,
			},
			"noti_senderId":   nil,
			"noti_receivedId": nil,
		},
	}

	// body, err := json.Marshal(messageNoti)
	// if err != nil {
	// 	return response.ErrorInsert, out, err
	// }

	err = helper.SendToRabbitMQ(messageNoti, consts.ExchangeNoti, consts.RoutingKeyNoti)
	if err != nil {
		return response.ErrorRabbitMQ, out, err
	}
	global.Logger.Info("Send to RabbitMQ success", zap.Any("message", messageNoti))

	return 201, newEvent, nil
}
func (evt *sEvent) EditEvent(ctx context.Context, id string, in *model.UpdateEventParams) (resultCode int, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, err
	}
	event, err := evt.r.GetEventById(ctx, uuidID)
	if err != nil {
		return response.ErrorDataNotExists, err
	}

	err = evt.r.UpdateEvent(ctx, database.UpdateEventParams{
		ID:               event.ID,
		EventName:        in.Name,
		EventDescription: pgtype.Text{String: in.Description, Valid: true},
		EventImageUrl:    pgtype.Text{String: in.ImageUrl, Valid: true},
		EventStart:       pgtype.Timestamp{Time: in.Start, Valid: true},
		EventEnd:         pgtype.Timestamp{Time: in.End, Valid: true},
		EventActive:      pgtype.Bool{Bool: in.Active},
		UserID:           pgtype.UUID{Bytes: uuid.MustParse(in.UserId), Valid: true},
	})
	if err != nil {
		return response.ErrorUpdate, err
	}

	return 200, nil
}
func (evt *sEvent) GetEventById(ctx context.Context, id string) (resultCode int, out database.PreGoEvent, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, out, err
	}
	event, err := evt.r.GetEventById(ctx, uuidID)
	if err != nil {
		return response.ErrorDataNotExists, out, err
	}

	if !event.ID.Valid {
		return response.ErrorListFailed, out, nil
	}

	return 200, event, nil
}
func (evt *sEvent) DeleteEvent(ctx context.Context, id string) (resultCode int, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, err
	}
	event, err := evt.r.GetEventById(ctx, uuidID)
	if err != nil {
		return response.ErrorDataNotExists, err
	}

	if !event.ID.Valid {
		return response.ErrorListFailed, nil
	}

	err = evt.r.DeleteEvent(ctx, event.ID)
	if err != nil {
		return response.ErrorDelete, err
	}

	return 200, nil
}

func (evt *sEvent) EventsLike(ctx context.Context, id string, userId string) (resultCode int, err error) {
	uuidID, err := ParseUUID(id)
	if err != nil {
		return response.ErrParseUUID, err
	}

	event, err := evt.r.GetEventById(ctx, uuidID)
	if err != nil {
		return response.ErrorDataNotExists, err
	}

	_, err = evt.r.CreateEventLike(ctx, database.CreateEventLikeParams{
		EventID: event.ID,
		UserID:  pgtype.UUID{Bytes: uuid.MustParse(userId), Valid: true},
	})
	if err != nil {
		return response.ErrorInsert, err
	}
	return 200, nil
}
func (evt *sEvent) EventsUnLike(ctx context.Context, eventID string, userID string) (resultCode int, err error) {
	parsedEventID, err := ParseUUID(eventID)
	if err != nil {
		return response.ErrParseUUID, err
	}

	event, err := evt.r.GetEventById(ctx, parsedEventID)
	if err != nil {
		return response.ErrorDataNotExists, err
	}

	parsedUserID := pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true}
	err = evt.r.DeleteEventLike(ctx, database.DeleteEventLikeParams{
		EventID: event.ID,
		UserID:  parsedUserID,
	})
	if err != nil {
		return response.ErrorDelete, err
	}
	return 200, nil
}

func (evt *sEvent) IsLiked(ctx context.Context, userId string) (resultCode int, out []database.GetEventsUserLikeRow, err error) {
	parsedUserId, err := ParseUUID(userId)
	if err != nil {
		return response.ErrParseUUID, out, err
	}

	eventUser, err := evt.r.GetEventsUserLike(ctx, parsedUserId)
	if err != nil {
		return response.ErrorListFailed, out, err
	}
	return 200, eventUser, nil
}

func ParseUUID(str string) (pgtype.UUID, error) {
	var uuidConvert pgtype.UUID
	err := uuidConvert.Scan(str)
	return uuidConvert, err
}
