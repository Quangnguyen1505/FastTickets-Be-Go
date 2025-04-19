package imple

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/model"
	"github.com/ntquang/ecommerce/response"
)

type sMenuFunction struct {
	r *database.Queries
}

func NewMenuFunctionImpl(r *database.Queries) *sMenuFunction {
	return &sMenuFunction{
		r: r,
	}
}

// implement
func (menuFunc *sMenuFunction) GetAllMenuFunctionsActive(ctx context.Context) (resultCode int, out []database.PreGoMenuFunction, err error) {
	menuFunction, err := menuFunc.r.GetAllMenuFunctionActive(ctx)
	if err != nil {
		return response.ErrorDataNotExists, nil, err
	}

	if len(menuFunction) == 0 {
		return response.ErrorDataNotExists, nil, err
	}

	return 200, menuFunction, nil
}

func (menuFunc *sMenuFunction) GetAllMenuFunctions(ctx context.Context) (resultCode int, out []database.PreGoMenuFunction, err error) {
	menuFunction, err := menuFunc.r.GetAllMenuFunctions(ctx)
	if err != nil {
		return response.ErrorDataNotExists, nil, err
	}

	if len(menuFunction) == 0 {
		return response.ErrorDataNotExists, nil, err
	}

	return 200, menuFunction, nil
}

func (menuFunc *sMenuFunction) NewMenuFunctions(ctx context.Context, in *model.NewOrUpdateMenuFunctionParams) (resultCode int, out database.PreGoMenuFunction, err error) {
	menuFunction, err := menuFunc.r.GetMenuFunctionByName(ctx, in.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return response.ErrorListFailed, out, err
	}
	if menuFunction.ID != "" {
		return response.ErrorListFailed, out, err
	}

	newMenuFunction, err := menuFunc.r.AddNewMenuFunction(ctx, database.AddNewMenuFunctionParams{
		Name:        in.Name,
		Description: pgtype.Text{String: in.Description, Valid: false},
		Url:         pgtype.Text{String: in.Url, Valid: true},
		Active:      pgtype.Bool{Bool: in.Active, Valid: true},
	})
	if err != nil {
		return response.ErrorInsert, out, nil
	}
	return 201, newMenuFunction, nil
}

func (menuFunc *sMenuFunction) EditMenuFunctionsById(ctx context.Context, id string, in *model.NewOrUpdateMenuFunctionParams) (resultCode int, out database.PreGoMenuFunction, err error) {
	menuFunction, err := menuFunc.r.GetMenuFunctionById(ctx, id)
	if err != nil {
		return response.ErrCodeOtpNotExists, out, err
	}

	editMenuFunc, err := menuFunc.r.EditMenuFunction(ctx, database.EditMenuFunctionParams{
		ID:          menuFunction.ID,
		Name:        in.Name,
		Description: pgtype.Text{String: in.Description, Valid: false},
		Url:         pgtype.Text{String: in.Url, Valid: false},
		Active:      pgtype.Bool{Bool: in.Active, Valid: true},
	})
	if err != nil {
		return response.ErrorInsert, out, err
	}

	return 200, editMenuFunc, err
}

func (menuFunc *sMenuFunction) GetMenuFunctionsById(ctx context.Context, id string) (resultCode int, out database.PreGoMenuFunction, err error) {
	menuFunction, err := menuFunc.r.GetMenuFunctionById(ctx, id)
	if err != nil {
		return response.ErrorDataNotExists, out, err
	}

	if menuFunction.ID == "" {
		return response.ErrorListFailed, out, err
	}

	return 200, menuFunction, nil
}

func (menuFunc *sMenuFunction) DeleteMenuFunctionsById(ctx context.Context, id string) (resultCode int, name string, err error) {
	menuFunction, err := menuFunc.r.GetMenuFunctionById(ctx, id)
	if err != nil {
		return response.ErrorDataNotExists, "", err
	}

	if menuFunction.ID == "" {
		return response.ErrorListFailed, "", err
	}

	menufuncDel, err := menuFunc.r.RemoveMenuFunction(ctx, menuFunction.ID)
	if err != nil {
		return response.ErrorDelete, "", err
	}

	return 200, menufuncDel.Name, nil
}
