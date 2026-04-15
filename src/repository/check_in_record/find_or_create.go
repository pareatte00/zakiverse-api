package check_in_record

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOrCreate(ctx context.Context, accountId string, planId string) (model.CheckInRecord, error) {
	var dest model.CheckInRecord

	// Try insert, on conflict do nothing
	insertStmt := CheckInRecord.INSERT(
		CheckInRecord.AccountID,
		CheckInRecord.PlanID,
	).VALUES(
		accountId,
		planId,
	).ON_CONFLICT(CheckInRecord.AccountID, CheckInRecord.PlanID).DO_NOTHING()

	_, _ = insertStmt.ExecContext(ctx, r.db)

	// Then select
	selectStmt := postgres.SELECT(CheckInRecord.AllColumns).
		FROM(CheckInRecord).
		WHERE(
			CheckInRecord.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()).
				AND(CheckInRecord.PlanID.EQ(postgres.CAST(postgres.String(planId)).AS_UUID())),
		)

	err := selectStmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}

func (r *Repository) FindByAccountId(ctx context.Context, accountId string) ([]model.CheckInRecord, error) {
	var dest []model.CheckInRecord

	stmt := postgres.SELECT(CheckInRecord.AllColumns).
		FROM(CheckInRecord).
		WHERE(CheckInRecord.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
