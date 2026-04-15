package account_balance

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) DeductCoins(ctx context.Context, accountId string, amount int) (model.AccountBalance, error) {
	var dest model.AccountBalance

	stmt := AccountBalance.UPDATE(AccountBalance.Coin, AccountBalance.UpdatedAt).
		SET(
			AccountBalance.Coin.SUB(postgres.Int32(int32(amount))),
			postgres.NOW(),
		).
		WHERE(
			AccountBalance.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()).
				AND(AccountBalance.Coin.GT_EQ(postgres.Int32(int32(amount)))),
		).
		RETURNING(AccountBalance.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
