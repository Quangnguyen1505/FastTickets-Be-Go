package imple

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ntquang/ecommerce/response"

	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
)

type sEvent struct {
	r *database.Queries
}

func NewEventImpl(r *database.Queries) *sEvent {
	return &sEvent{
		r: r,
	}
}

// implement
func (evt *sEvent) GetAllEventsActive(ctx context.Context, query model.EventQuery) (resultCode int, out []database.PreGoEvent, err error) {
	limit, page := query.Limit, query.Page
	offset := (page - 1) * limit

	params := database.GetAllActiveEventsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	foundEvent, err := evt.r.GetAllActiveEvents(ctx, params)
	if err != nil {
		return response.ErrorListFailed, nil, err
	}

	if len(foundEvent) == 0 {
		return response.ErrorDataNotExists, nil, nil
	}

	return 200, foundEvent, nil
}

func (evt *sEvent) NewEvent(ctx context.Context, userId string, in *model.AddNewEventParams) (resultCode int, out database.PreGoEvent, err error) {
	fmt.Println("evenet name and des ", in.Name, in.Description)
	event, err := evt.r.GetEventByName(ctx, in.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return response.ErrorListFailed, out, err
	}
	if event.ID.Valid {
		return response.ErrorListFailed, out, err
	}
	fmt.Println("evenet name and des ", in.Name, in.Description)
	newEvent, err := evt.r.AddNewEvent(ctx, database.AddNewEventParams{
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

func ParseUUID(str string) (pgtype.UUID, error) {
	var uuidConvert pgtype.UUID
	err := uuidConvert.Scan(str)
	return uuidConvert, err
}
