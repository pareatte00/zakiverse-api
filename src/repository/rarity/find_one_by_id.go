package rarity

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOneById(ctx context.Context, id string) (model.Rarity, error) {
	var dest model.Rarity

	stmt := postgres.SELECT(Rarity.AllColumns).
		FROM(Rarity).
		WHERE(Rarity.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
