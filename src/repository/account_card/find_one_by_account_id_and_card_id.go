package account_card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOneByAccountIdAndCardId(ctx context.Context, accountId string, cardId string) (model.AccountCard, error) {
	var dest model.AccountCard

	stmt := postgres.SELECT(AccountCard.AllColumns).
		FROM(AccountCard).
		WHERE(
			AccountCard.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()).
				AND(AccountCard.CardID.EQ(postgres.CAST(postgres.String(cardId)).AS_UUID())),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
