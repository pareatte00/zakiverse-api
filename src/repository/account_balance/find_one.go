package account_balance

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOne(ctx context.Context, accountId string) (model.AccountBalance, error) {
	var dest model.AccountBalance

	stmt := postgres.SELECT(AccountBalance.AllColumns).
		FROM(AccountBalance).
		WHERE(AccountBalance.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
