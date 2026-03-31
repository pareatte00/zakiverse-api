package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllParam struct {
	ActiveOnly bool
	Limit      int64
	Offset     int64
}

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]model.Pack, error) {
	var dest []model.Pack

	condition := postgres.Bool(true)
	if param.ActiveOnly {
		condition = condition.AND(Pack.IsActive.EQ(postgres.Bool(true)))
	}

	stmt := postgres.SELECT(Pack.AllColumns).
		FROM(Pack).
		WHERE(condition).
		ORDER_BY(Pack.CreatedAt.DESC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
