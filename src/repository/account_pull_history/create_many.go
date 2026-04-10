package account_pull_history

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateManyParam struct {
	AccountId  string
	PackId     string
	CardId     string
	Rarity     string
	IsPity     bool
	IsFeatured bool
	IsNew      bool
}

func (r *Repository) CreateMany(ctx context.Context, params []CreateManyParam) error {
	if len(params) == 0 {
		return nil
	}

	stmt := AccountPullHistory.INSERT(
		AccountPullHistory.AccountID,
		AccountPullHistory.PackID,
		AccountPullHistory.CardID,
		AccountPullHistory.Rarity,
		AccountPullHistory.IsPity,
		AccountPullHistory.IsFeatured,
		AccountPullHistory.IsNew,
	)

	for _, p := range params {
		stmt = stmt.VALUES(
			postgres.CAST(postgres.String(p.AccountId)).AS_UUID(),
			postgres.CAST(postgres.String(p.PackId)).AS_UUID(),
			postgres.CAST(postgres.String(p.CardId)).AS_UUID(),
			postgres.NewEnumValue(p.Rarity),
			p.IsPity,
			p.IsFeatured,
			p.IsNew,
		)
	}

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
