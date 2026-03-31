package account_pack_pity

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindAllByAccountAndPack(ctx context.Context, accountId string, packId string) ([]model.AccountPackPity, error) {
	var dest []model.AccountPackPity

	stmt := postgres.SELECT(AccountPackPity.AllColumns).
		FROM(AccountPackPity).
		WHERE(
			AccountPackPity.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()).
				AND(AccountPackPity.PackID.EQ(postgres.CAST(postgres.String(packId)).AS_UUID())),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
