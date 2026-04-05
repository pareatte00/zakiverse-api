package pack_pool

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Name        string
	Description *string
	ActiveCount int32
	RotationDay int32
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.PackPool, error) {
	var dest model.PackPool

	stmt := PackPool.INSERT(
		PackPool.Name,
		PackPool.Description,
		PackPool.ActiveCount,
		PackPool.RotationDay,
	).VALUES(
		param.Name,
		param.Description,
		param.ActiveCount,
		param.RotationDay,
	).RETURNING(PackPool.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
