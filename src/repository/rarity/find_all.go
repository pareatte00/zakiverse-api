package rarity

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindAll(ctx context.Context) ([]model.Rarity, error) {
	var dest []model.Rarity

	stmt := postgres.SELECT(Rarity.AllColumns).
		FROM(Rarity).
		ORDER_BY(Rarity.Name.ASC())

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
