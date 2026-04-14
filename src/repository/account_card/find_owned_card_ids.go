package account_card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOwnedCardIds(ctx context.Context, accountId string, cardIds []uuid.UUID) (map[uuid.UUID]bool, error) {
	if len(cardIds) == 0 {
		return map[uuid.UUID]bool{}, nil
	}

	expressions := make([]postgres.Expression, len(cardIds))
	for i, id := range cardIds {
		expressions[i] = postgres.CAST(postgres.String(id.String())).AS_UUID()
	}

	var dest []struct {
		CardId uuid.UUID `alias:"account_card.card_id"`
	}

	stmt := postgres.SELECT(AccountCard.CardID).
		FROM(AccountCard).
		WHERE(
			AccountCard.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()).
				AND(AccountCard.CardID.IN(expressions...)),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	result := make(map[uuid.UUID]bool, len(dest))
	for _, d := range dest {
		result[d.CardId] = true
	}

	return result, nil
}
