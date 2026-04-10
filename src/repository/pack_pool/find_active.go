package pack_pool

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindActive(ctx context.Context) ([]model.PackPool, error) {
	var dest []model.PackPool

	stmt := postgres.SELECT(PackPool.AllColumns).
		FROM(PackPool).
		WHERE(PackPool.IsActive.EQ(postgres.Bool(true))).
		ORDER_BY(PackPool.SortOrder.ASC(), PackPool.CreatedAt.DESC())

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
