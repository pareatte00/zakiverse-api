package account_card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) DeleteOne(ctx context.Context, accountId string, cardId string) error {
	stmt := AccountCard.DELETE().
		WHERE(
			AccountCard.AccountID.EQ(postgres.String(accountId)).
				AND(AccountCard.CardID.EQ(postgres.String(cardId))),
		)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
