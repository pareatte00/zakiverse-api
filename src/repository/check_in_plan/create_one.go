package check_in_plan

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Code        string
	Name        string
	Description *string
	Type        string
	Interval    int32
	MaxClaims   int32
	Rewards     string
	ResetPolicy string
	IsActive    bool
	StartsAt    interface{}
	EndsAt      interface{}
	SortOrder   int32
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.CheckInPlan, error) {
	var dest model.CheckInPlan

	stmt := CheckInPlan.INSERT(
		CheckInPlan.Code,
		CheckInPlan.Name,
		CheckInPlan.Description,
		CheckInPlan.Type,
		CheckInPlan.Interval,
		CheckInPlan.MaxClaims,
		CheckInPlan.Rewards,
		CheckInPlan.ResetPolicy,
		CheckInPlan.IsActive,
		CheckInPlan.StartsAt,
		CheckInPlan.EndsAt,
		CheckInPlan.SortOrder,
	).VALUES(
		param.Code,
		param.Name,
		param.Description,
		param.Type,
		param.Interval,
		param.MaxClaims,
		param.Rewards,
		param.ResetPolicy,
		param.IsActive,
		param.StartsAt,
		param.EndsAt,
		param.SortOrder,
	).RETURNING(CheckInPlan.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
