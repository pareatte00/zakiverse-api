package check_in_plan

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllParam struct {
	Limit  int64
	Offset int64
}

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]model.CheckInPlan, error) {
	var dest []model.CheckInPlan

	stmt := postgres.SELECT(CheckInPlan.AllColumns).
		FROM(CheckInPlan).
		ORDER_BY(CheckInPlan.SortOrder.ASC(), CheckInPlan.CreatedAt.DESC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}

func (r *Repository) FindActive(ctx context.Context) ([]model.CheckInPlan, error) {
	var dest []model.CheckInPlan

	stmt := postgres.SELECT(CheckInPlan.AllColumns).
		FROM(CheckInPlan).
		WHERE(
			CheckInPlan.IsActive.EQ(postgres.Bool(true)).AND(
				CheckInPlan.StartsAt.IS_NULL().OR(CheckInPlan.StartsAt.LT_EQ(postgres.TimestampzT(time.Now()))),
			).AND(
				CheckInPlan.EndsAt.IS_NULL().OR(CheckInPlan.EndsAt.GT(postgres.TimestampzT(time.Now()))),
			),
		).
		ORDER_BY(CheckInPlan.SortOrder.ASC())

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
