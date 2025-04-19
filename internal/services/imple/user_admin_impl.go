package imple

import (
	"context"

	"github.com/ntquang/ecommerce/internal/database"
)

type sUserAdmin struct {
	r *database.Queries
}

func NewUserAdminImpl(r *database.Queries) *sUserAdmin {
	return &sUserAdmin{
		r: r,
	}
}

func (su *sUserAdmin) RemoveUser(ctx context.Context) (result string) {
	result = "xoa thanh cong"
	return result
}

func (su *sUserAdmin) FindOneUser(ctx context.Context) error {
	return nil
}
