package account_card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllByAccountIdParam struct {
	AccountId string
	Limit     int64
	Offset    int64
}

func (r *Repository) FindAllByAccountId(ctx context.Context, param FindAllByAccountIdParam) ([]model.AccountCard, error) {
	var dest []model.AccountCard

	stmt := postgres.SELECT(AccountCard.AllColumns).
		FROM(AccountCard).
		WHERE(AccountCard.AccountID.EQ(postgres.CAST(postgres.String(param.AccountId)).AS_UUID())).
		ORDER_BY(AccountCard.ObtainedAt.DESC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
