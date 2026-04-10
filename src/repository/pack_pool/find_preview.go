package pack_pool

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

// FindPreview returns pools that are not yet active but within their preview window.
// A pool is in preview when: open_at is set, is_active=false, preview_days>0,
// and open_at - preview_days <= now < open_at.
func (r *Repository) FindPreview(ctx context.Context, now time.Time) ([]model.PackPool, error) {
	var dest []model.PackPool

	stmt := postgres.SELECT(PackPool.AllColumns).
		FROM(PackPool).
		WHERE(
			PackPool.IsActive.EQ(postgres.Bool(false)).
				AND(PackPool.PreviewDays.GT(postgres.Int(0))).
				AND(PackPool.OpenAt.IS_NOT_NULL()).
				AND(PackPool.OpenAt.GT(postgres.TimestampzT(now))).
				AND(postgres.RawTimestampz("pack_pool.open_at - (pack_pool.preview_days || ' days')::interval").LT_EQ(postgres.TimestampzT(now))),
		).
		ORDER_BY(PackPool.SortOrder.ASC())

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
