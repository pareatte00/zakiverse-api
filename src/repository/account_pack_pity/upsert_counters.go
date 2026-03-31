package account_pack_pity

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type UpsertParam struct {
	AccountId string
	PackId    string
	Rarity    string
	Counter   int32
}

func (r *Repository) UpsertCounters(ctx context.Context, params []UpsertParam) error {
	if len(params) == 0 {
		return nil
	}

	stmt := AccountPackPity.INSERT(
		AccountPackPity.AccountID,
		AccountPackPity.PackID,
		AccountPackPity.Rarity,
		AccountPackPity.Counter,
	)

	for _, p := range params {
		stmt = stmt.VALUES(p.AccountId, p.PackId, p.Rarity, p.Counter)
	}

	stmt = stmt.ON_CONFLICT(AccountPackPity.AccountID, AccountPackPity.PackID, AccountPackPity.Rarity).
		DO_UPDATE(postgres.SET(
			AccountPackPity.Counter.SET(AccountPackPity.EXCLUDED.Counter),
		))

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
