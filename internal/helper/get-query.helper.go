package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/internal/model"
)

func ParseEventQuery(ctx *gin.Context) model.EventQuery {
	return model.EventQuery{
		Limit: GetQueryInt(ctx, "limit", 10),
		Page:  GetQueryInt(ctx, "page", 1),
	}
}

func ParseContactMessageQuery(ctx *gin.Context) model.ContactMessageQuery {
	status := int16(GetQueryInt(ctx, "status", 0))
	var statusPtr *int16
	if status != 0 {
		statusPtr = &status
	}

	return model.ContactMessageQuery{
		Status: statusPtr,
		Limit:  GetQueryInt(ctx, "limit", 10),
		Page:   GetQueryInt(ctx, "page", 1),
	}
}

func GetQueryInt(ctx *gin.Context, key string, defaultVal int) int {
	valStr := ctx.Query(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
