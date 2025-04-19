package menufunction

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/internal/services"
	"github.com/ntquang/ecommerce/response"
	"go.uber.org/zap"
)

var MenuFunc = new(cMenuFunction)

type cMenuFunction struct{}

// @Summary      Get all active menu functions
// @Description  Get all menu functions that are currently active
// @Tags         menu function
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /menu-functions/active [get]
func (cMenuFunc *cMenuFunction) GetAllMenuFunctionsActive(ctx *gin.Context) {
	statusCode, menuFunc, err := services.MenuFunction().GetAllMenuFunctionsActive(ctx)
	if err != nil {
		global.Logger.Error("Error get active menu functions", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, menuFunc, "Get all active menu functions successfully!")
}

// @Summary      Get all menu functions
// @Description  Get all menu functions
// @Tags         menu function
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /menu-functions [get]
func (cMenuFunc *cMenuFunction) GetAllMenuFunctions(ctx *gin.Context) {
	statusCode, menuFunc, err := services.MenuFunction().GetAllMenuFunctions(ctx)
	if err != nil {
		global.Logger.Error("Error get menu functions", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, menuFunc, "Get all menu functions successfully!")
}

// @Summary      Create new menu function
// @Description  Create a new menu function
// @Tags         menu function
// @Accept       json
// @Produce      json
// @Param        payload body model.NewOrUpdateMenuFunctionParams true "payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrResponse
// @Router       /menu-functions [post]
func (cMenuFunc *cMenuFunction) NewMenuFunctions(ctx *gin.Context) {
	var params model.NewOrUpdateMenuFunctionParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		global.Logger.Error("Error parsing menu function params", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}
	statusCode, metadata, err := services.MenuFunction().NewMenuFunctions(ctx, &params)
	if err != nil {
		global.Logger.Error("Error creating menu function", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, metadata, "New menu function created successfully!")
}

// @Summary      Update menu function
// @Description  Update an existing menu function
// @Tags         menu function
// @Accept       json
// @Produce      json
// @Param        id      path      int                        true  "Menu Function ID"
// @Param        payload body      model.NewOrUpdateMenuFunctionParams true  "Payload"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /menu-functions/{id} [put]
func (cMenuFunc *cMenuFunction) EditMenuFunctionsById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var params model.NewOrUpdateMenuFunctionParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		global.Logger.Error("Error parsing menu function params", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "", fmt.Errorf(err.Error()))
		return
	}
	statusCode, metadata, err := services.MenuFunction().EditMenuFunctionsById(ctx, idStr, &params)
	if err != nil {
		global.Logger.Error("Error updating menu function", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, metadata, "Menu function updated successfully!")
}

// @Summary      Get menu function by ID
// @Description  Retrieve a specific menu function by its ID
// @Tags         menu function
// @Accept       json
// @Produce      json
// @Param        id      path      int                        true  "Menu Function ID"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /menu-functions/{id} [get]
func (cMenuFunc *cMenuFunction) GetMenuFunctionsById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	statusCode, metadata, err := services.MenuFunction().GetMenuFunctionsById(ctx, idStr)
	if err != nil {
		global.Logger.Error("Error retrieving menu function", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, metadata, "Get menu function by ID successfully!")
}

// @Summary      Delete menu function
// @Description  Delete a menu function by its ID
// @Tags         menu function
// @Accept       json
// @Produce      json
// @Param        id      path      int                        true  "Menu Function ID"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrResponse
// @Router       /menu-function/{id} [delete]
func (cMenuFunc *cMenuFunction) DeleteMenuFunctionsById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	statusCode, name, err := services.MenuFunction().DeleteMenuFunctionsById(ctx, idStr)
	if err != nil {
		global.Logger.Error("Error deleting menu function", zap.Error(err))
		response.ErrorResponse(ctx, statusCode, "", err)
		return
	}
	response.SuccessResponse(ctx, statusCode, name, "Menu function deleted successfully!")
}
