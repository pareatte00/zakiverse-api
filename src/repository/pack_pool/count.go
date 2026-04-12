package pack_pool

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CountParam struct {
	BannerType string
	ActiveOnly bool
}

func (r *Repository) Count(ctx context.Context, param CountParam) (int64, error) {
	var dest struct {
		Count int64
	}

	condition := postgres.Bool(true)

	if param.BannerType != "" {
		condition = condition.AND(PackPool.BannerType.EQ(postgres.NewEnumValue(param.BannerType)))
	}

	if param.ActiveOnly {
		condition = condition.AND(PackPool.IsActive.EQ(postgres.Bool(true)))
	}

	stmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("count"),
	).FROM(PackPool).
		WHERE(condition)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}
