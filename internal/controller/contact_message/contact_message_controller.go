package contactmessage

import (
	"fmt"

	"github.com/ntquang/ecommerce/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/internal/services"
	"github.com/ntquang/ecommerce/response"
	"go.uber.org/zap"
)

var ContactMessage = new(cContactMessage)

type cContactMessage struct{}

// ContactMessage Documentation
// @Summary      Create new contact message
// @Description  When user sends a new contact message
// @Tags         contactmessage
// @Accept       json
// @Produce      json
// @Param        payload body model.AddNewContactMessageParams true "Payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /contact-messages [post]
func (cm *cContactMessage) NewContactMessage(ctx *gin.Context) {
	var params model.AddNewContactMessageParams
	if err := ctx.ShouldBind(&params); err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params event", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	statusCode, metadata, err := services.ContactMessage().NewContactMessage(ctx, &params)
	if err != nil {
		global.Logger.Error("Error create event", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, statusCode, metadata, "New Contact Messages Created Successfully!")
}

// ContactMessage Documentation
// @Summary      Get all contact messages by status
// @Description  When admin wants to get all contact messages filtered by status
// @Tags         contactmessage
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /contact-messages [get]
func (cm *cContactMessage) GetAllContactMessageStatus(ctx *gin.Context) {
	query := helper.ParseContactMessageQuery(ctx)
	statusCode, metadata, err := services.ContactMessage().GetAllContactMessageByStatus(ctx, query)
	if err != nil {
		global.Logger.Error("Error get contact message", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, metadata, "Get All Contact Message By Status Successfully!")
}

// ContactMessage Documentation
// @Summary      Update contact message
// @Description  When admin wants to update status of a contact message
// @Tags         contactmessage
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Param        id      path      string                          true  "Contact Message ID"
// @Param        payload body      model.UpdateContactMessageParams true  "Payload"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /contact-messages/{id} [put]
func (cm *cContactMessage) EditContactMessage(ctx *gin.Context) {
	var params map[string]interface{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		fmt.Println("Lỗi bind JSON:", err)
		global.Logger.Error("Error parse params contact message", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	status, ok := params["status"].(float64)
	if !ok {
		global.Logger.Error("Missing or invalid status parameter")
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf("Invalid or missing status"))
		return
	}

	contactId := ctx.Param("id")
	statusCode, err := services.ContactMessage().EditStatusContactMessage(ctx, contactId, int16(status))
	if err != nil {
		global.Logger.Error("Error update contact message", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, 200, nil, "Update contact message Successfully!")
}

// ContactMessage Documentation
// @Summary      Get contact message by ID
// @Description  When admin wants to get a contact message by its ID
// @Tags         contactmessage
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Contact Message ID"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /contact-messages/{id} [get]
func (cm *cContactMessage) GetContactMessageById(ctx *gin.Context) {
	contactId := ctx.Param("id")
	statusCode, metadata, err := services.ContactMessage().GetContactMessageById(ctx, contactId)
	if err != nil {
		global.Logger.Error("Error get contact message by ID", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", fmt.Errorf(err.Error()))
		return
	}
	response.SuccessResponse(ctx, 200, metadata, "Get Event By ID Successfully!")
}

// ContactMessage Documentation
// @Summary      Delete contact message
// @Description  When admin wants to delete a contact message
// @Tags         contactmessage
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Param        id      path      string  true  "Contact Message ID"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /contact-messages/{id} [delete]
func (cm *cContactMessage) DeleteContactMessage(ctx *gin.Context) {
	contactId := ctx.Param("id")
	statusCode, err := services.ContactMessage().DeleteContactMessage(ctx, contactId)
	if err != nil {
		global.Logger.Error("Error delete event", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, 200, nil, "Delete Contact Message Successfully!")
}

// ContactMessage Documentation
// @Summary      Send email to customer
// @Description  When admin wants to send an email to a customer
// @Tags         contactmessage
// @Accept       json
// @Produce      json
// @Param        authorization header string true "authorization token"
// @Param        x-client-id header string true "x-client-id user"
// @Param        payload body      model.ResponseCustomer true  "Payload"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /contact-messages/customer [post]
func (c *cContactMessage) SendEmailToCustomer(ctx *gin.Context) {
	var params model.ResponseCustomer
	if err := ctx.ShouldBindJSON(&params); err != nil {
		global.Logger.Error("Error parse params contact message", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}

	statusCode, err := services.ContactMessage().SendEmailToCustomer(ctx, &params)
	if err != nil {
		global.Logger.Error("Error send email to customer", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", fmt.Errorf(err.Error()))
		return
	}

	response.SuccessResponse(ctx, 200, nil, "Send Email To Customer Successfully!")
}
