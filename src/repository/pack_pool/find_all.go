package pack_pool

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllParam struct {
	BannerType string
	ActiveOnly bool
	Limit      int64
	Offset     int64
}

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]model.PackPool, error) {
	var dest []model.PackPool

	condition := postgres.Bool(true)

	if param.BannerType != "" {
		condition = condition.AND(PackPool.BannerType.EQ(postgres.NewEnumValue(param.BannerType)))
	}

	if param.ActiveOnly {
		condition = condition.AND(PackPool.IsActive.EQ(postgres.Bool(true)))
	}

	stmt := postgres.SELECT(PackPool.AllColumns).
		FROM(PackPool).
		WHERE(condition).
		ORDER_BY(PackPool.SortOrder.ASC(), PackPool.CreatedAt.DESC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
