package rarity

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type UpdateOneByIdParam struct {
	Name   string
	Config string
}

func (r *Repository) UpdateOneById(ctx context.Context, id string, param UpdateOneByIdParam) (model.Rarity, error) {
	var dest model.Rarity

	stmt := Rarity.UPDATE(
		Rarity.Name,
		Rarity.Config,
	).SET(
		param.Name,
		param.Config,
	).WHERE(
		Rarity.ID.EQ(postgres.String(id)),
	).RETURNING(Rarity.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
