package check_in_plan

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) UpdateOneById(ctx context.Context, id string, updates map[string]any) (model.CheckInPlan, error) {
	var dest model.CheckInPlan

	columnMap := map[string]postgres.Column{
		"code":         CheckInPlan.Code,
		"name":         CheckInPlan.Name,
		"description":  CheckInPlan.Description,
		"type":         CheckInPlan.Type,
		"interval":     CheckInPlan.Interval,
		"max_claims":   CheckInPlan.MaxClaims,
		"rewards":      CheckInPlan.Rewards,
		"reset_policy": CheckInPlan.ResetPolicy,
		"is_active":    CheckInPlan.IsActive,
		"starts_at":    CheckInPlan.StartsAt,
		"ends_at":      CheckInPlan.EndsAt,
		"sort_order":   CheckInPlan.SortOrder,
	}

	var f postgres.ColumnList
	var vs []any
	for k, v := range updates {
		f = append(f, columnMap[k])
		vs = append(vs, v)
	}

	stmt := CheckInPlan.UPDATE(f).SET(vs[0], vs[1:]...).WHERE(
		CheckInPlan.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()),
	).RETURNING(CheckInPlan.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
