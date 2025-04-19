package services

import (
	"context"

	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
)

type (
	IMenuFunction interface {
		GetAllMenuFunctionsActive(ctx context.Context) (resultCode int, out []database.PreGoMenuFunction, err error)
		GetAllMenuFunctions(ctx context.Context) (resultCode int, out []database.PreGoMenuFunction, err error)
		NewMenuFunctions(ctx context.Context, in *model.NewOrUpdateMenuFunctionParams) (resultCode int, out database.PreGoMenuFunction, err error)
		EditMenuFunctionsById(ctx context.Context, id string, in *model.NewOrUpdateMenuFunctionParams) (resultCode int, out database.PreGoMenuFunction, err error)
		GetMenuFunctionsById(ctx context.Context, id string) (resultCode int, out database.PreGoMenuFunction, err error)
		DeleteMenuFunctionsById(ctx context.Context, id string) (resultCode int, name string, err error)
	}
)

var (
	localMenuFunction IMenuFunction
)

func MenuFunction() IMenuFunction {
	if localMenuFunction == nil {
		panic("implement localUserAdmin not found for interface IUserAdmin")
	}

	return localMenuFunction
}

func InitMenuFunction(i IMenuFunction) {
	localMenuFunction = i
}
