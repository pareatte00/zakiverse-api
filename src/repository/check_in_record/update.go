package check_in_record

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type UpdateParam struct {
	AccountId   string
	PlanId      string
	ClaimCount  int32
	Streak      int32
	LastClaimed time.Time
	ResetAt     *time.Time
}

func (r *Repository) Update(ctx context.Context, param UpdateParam) (model.CheckInRecord, error) {
	var dest model.CheckInRecord

	cols := postgres.ColumnList{
		CheckInRecord.ClaimCount,
		CheckInRecord.Streak,
		CheckInRecord.LastClaimed,
		CheckInRecord.UpdatedAt,
	}

	vals := []interface{}{
		param.ClaimCount,
		param.Streak,
		param.LastClaimed,
		postgres.NOW(),
	}

	if param.ResetAt != nil {
		cols = append(cols, CheckInRecord.ResetAt)
		vals = append(vals, *param.ResetAt)
	}

	stmt := CheckInRecord.UPDATE(cols).
		SET(vals[0], vals[1:]...).
		WHERE(
			CheckInRecord.AccountID.EQ(postgres.CAST(postgres.String(param.AccountId)).AS_UUID()).
				AND(CheckInRecord.PlanID.EQ(postgres.CAST(postgres.String(param.PlanId)).AS_UUID())),
		).
		RETURNING(CheckInRecord.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
