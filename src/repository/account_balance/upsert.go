package account_balance

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) Upsert(ctx context.Context, accountId string, coin int32) (model.AccountBalance, error) {
	var dest model.AccountBalance

	stmt := AccountBalance.INSERT(
		AccountBalance.AccountID,
		AccountBalance.Coin,
	).VALUES(
		accountId,
		coin,
	).ON_CONFLICT(AccountBalance.AccountID).DO_UPDATE(
		postgres.SET(
			AccountBalance.Coin.SET(AccountBalance.EXCLUDED.Coin),
			AccountBalance.UpdatedAt.SET(postgres.NOW()),
		),
	).RETURNING(AccountBalance.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
