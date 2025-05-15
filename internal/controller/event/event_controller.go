package event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/helper"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/internal/services"
	"github.com/ntquang/ecommerce/internal/utils/context"
	"github.com/ntquang/ecommerce/response"
	"go.uber.org/zap"
)

var Event = new(cEvent)

type cEvent struct{}

// Event Documentation
// @Summary      Get all events
// @Description  When admin wants to get all event content
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /events [get]
func (h *cEvent) GetAllEvents(ctx *gin.Context) {
	query := helper.ParseEventQuery(ctx)
	statusCode, metadata, err := services.Event().GetAllEventsActive(ctx, query)
	if err != nil {
		global.Logger.Error("Error get event", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, metadata, "Get All Events Successfully!")
}

// Event Documentation
// @Summary      Create new event
// @Description  When admin wants to add a new event
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Param        payload body model.AddNewEventParams true "payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /events [post]
func (h *cEvent) NewEvent(ctx *gin.Context) {
	var params model.AddNewEventParams
	if err := ctx.ShouldBind(&params); err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params event", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	//get userId from uuid (token)
	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrTwoFactorAuthSetUpFailed, "Missing get UUID", err)
	}

	imageURL, _ := ctx.Get("image_url")
	params.ImageUrl = imageURL.(string)

	statusCode, metadata, err := services.Event().NewEvent(ctx, userId, &params)
	if err != nil {
		global.Logger.Error("Error create event", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", fmt.Errorf(err.Error()))
		return
	}

	if statusCode == response.ErrorListFailed {
		response.ErrorResponse(ctx, response.ErrorInsert, "", fmt.Errorf("event already exists!"))
		return
	}
	response.SuccessResponse(ctx, statusCode, metadata, "New Event Created Successfully!")
}

// Event Documentation
// @Summary      Update event
// @Description  When admin wants to update an event
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Event ID"
// @Param        payload body      model.UpdateEventParams true  "Payload"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /events/{id} [put]
func (h *cEvent) EditEvent(ctx *gin.Context) {
	var params model.UpdateEventParams
	if err := ctx.ShouldBind(&params); err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params event", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	eventId := ctx.Param("id")
	imageURL, _ := ctx.Get("image_url")
	if imageURL != nil {
		params.ImageUrl = imageURL.(string)
	}

	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrTwoFactorAuthSetUpFailed, "Missing get UUID", err)
	}

	params.UserId = userId
	statusCode, err := services.Event().EditEvent(ctx, eventId, &params)
	if err != nil {
		global.Logger.Error("Error update event", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, 200, nil, "Update Event Successfully!")
}

// Event Documentation
// @Summary      Get event by ID
// @Description  When admin wants to get an event by ID
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Param        id      path      string  true  "Event ID"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /events/{id} [get]
func (h *cEvent) GetEventById(ctx *gin.Context) {
	eventId := ctx.Param("id")
	statusCode, metadata, err := services.Event().GetEventById(ctx, eventId)
	if err != nil {
		global.Logger.Error("Error get event by ID", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", fmt.Errorf(err.Error()))
		return
	}
	response.SuccessResponse(ctx, 200, metadata, "Get Event By ID Successfully!")
}

// Event Documentation
// @Summary      Delete event
// @Description  When admin wants to delete an event
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Param        id      path      string  true  "Event ID"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /events/{id} [delete]
func (h *cEvent) DeleteEvent(ctx *gin.Context) {
	eventId := ctx.Param("id")
	statusCode, err := services.Event().DeleteEvent(ctx, eventId)
	if err != nil {
		global.Logger.Error("Error delete event", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, 200, nil, "Delete Event Successfully!")
}

// Event Documentation
// @Summary      Like an event
// @Description  When a user likes an event
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Event ID"
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /events/{id}/like [post]
func (h *cEvent) EventLike(ctx *gin.Context) {
	eventId := ctx.Param("id")

	//get userId from uuid (token)
	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrTwoFactorAuthSetUpFailed, "Missing get UUID", err)
	}

	statusCode, err := services.Event().EventsLike(ctx, eventId, userId)
	if err != nil {
		global.Logger.Error("Error event like", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}

	response.SuccessResponse(ctx, 200, nil, "Event Like Successfully!")
}

// Event Documentation
// @Summary      Unlike an event
// @Description  When a user unlikes an event
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Event ID"
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /events/{id}/unlike [delete]
func (h *cEvent) EventUnLike(ctx *gin.Context) {
	eventId := ctx.Param("id")

	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrTwoFactorAuthSetUpFailed, "Missing get UUID", err)
	}

	statusCode, err := services.Event().EventsUnLike(ctx, eventId, userId)

	if err != nil {
		global.Logger.Error("Error event unlike", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}

	response.SuccessResponse(ctx, 200, nil, "Event Unlike Successfully!")
}

// Event Documentation
// @Summary      get events like by user
// @Description  get events like by user
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /events/users [get]
func (h *cEvent) EventsIsLike(ctx *gin.Context) {
	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrTwoFactorAuthSetUpFailed, "Missing get UUID", err)
	}

	statusCode, metadata, err := services.Event().IsLiked(ctx, userId)

	if err != nil {
		global.Logger.Error("Error event unlike", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}

	response.SuccessResponse(ctx, 200, metadata, "Get event like by user Successfully!")
}
