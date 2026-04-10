package account_pull_history

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindByAccountParam struct {
	AccountId string
	Limit     int64
	Offset    int64
}

func (r *Repository) FindByAccount(ctx context.Context, param FindByAccountParam) ([]model.AccountPullHistory, error) {
	var dest []model.AccountPullHistory

	stmt := postgres.SELECT(AccountPullHistory.AllColumns).
		FROM(AccountPullHistory).
		WHERE(
			AccountPullHistory.AccountID.EQ(postgres.CAST(postgres.String(param.AccountId)).AS_UUID()),
		).
		ORDER_BY(AccountPullHistory.PulledAt.DESC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
