package rarity

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Name   string
	Config string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.Rarity, error) {
	var dest model.Rarity

	stmt := Rarity.INSERT(
		Rarity.Name,
		Rarity.Config,
	).VALUES(
		param.Name,
		param.Config,
	).RETURNING(Rarity.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
