package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Code          string
	Name          string
	Description   *string
	Image         string
	NameImage     *string
	CardsPerPull  int32
	SortOrder     int32
	Config        string
	PoolId        string
	RotationOrder *int32
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.Pack, error) {
	var dest model.Pack

	stmt := Pack.INSERT(
		Pack.Code,
		Pack.Name,
		Pack.Description,
		Pack.Image,
		Pack.NameImage,
		Pack.CardsPerPull,
		Pack.SortOrder,
		Pack.Config,
		Pack.PoolID,
		Pack.RotationOrder,
	).VALUES(
		param.Code,
		param.Name,
		param.Description,
		param.Image,
		param.NameImage,
		param.CardsPerPull,
		param.SortOrder,
		param.Config,
		postgres.CAST(postgres.String(param.PoolId)).AS_UUID(),
		param.RotationOrder,
	).RETURNING(Pack.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
