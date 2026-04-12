package account_pull_history

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) CountByAccountAndPack(ctx context.Context, accountId, packId string) (int64, error) {
	var dest struct {
		Count int64
	}

	stmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("count"),
	).FROM(AccountPullHistory).
		WHERE(
			AccountPullHistory.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()).
				AND(AccountPullHistory.PackID.EQ(postgres.CAST(postgres.String(packId)).AS_UUID())),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}
