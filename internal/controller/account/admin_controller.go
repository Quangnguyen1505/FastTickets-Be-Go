package account

import (
	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/internal/services"
	"github.com/ntquang/ecommerce/response"
)

var UserAdmin = new(cUserAdmin)

type cUserAdmin struct{}

func (c *cUserAdmin) RemoveUser(ctx *gin.Context) {
	result := services.UserAdmin().RemoveUser(ctx)

	response.SuccessResponse(ctx, response.ErrCodeRemoveSuccess, result, "")
}
