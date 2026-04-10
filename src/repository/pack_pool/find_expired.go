package pack_pool

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindExpired(ctx context.Context, now time.Time) ([]model.PackPool, error) {
	var dest []model.PackPool

	stmt := postgres.SELECT(PackPool.AllColumns).
		FROM(PackPool).
		WHERE(
			PackPool.IsActive.EQ(postgres.Bool(true)).
				AND(PackPool.CloseAt.IS_NOT_NULL()).
				AND(PackPool.CloseAt.LT_EQ(postgres.TimestampzT(now))),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
