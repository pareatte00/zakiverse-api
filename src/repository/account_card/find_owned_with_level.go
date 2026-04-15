package account_card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type OwnedCardInfo struct {
	CardId uuid.UUID
	Level  int32
}

func (r *Repository) FindOwnedWithLevel(ctx context.Context, accountId string, cardIds []uuid.UUID) (map[uuid.UUID]OwnedCardInfo, error) {
	if len(cardIds) == 0 {
		return map[uuid.UUID]OwnedCardInfo{}, nil
	}

	expressions := make([]postgres.Expression, len(cardIds))
	for i, id := range cardIds {
		expressions[i] = postgres.CAST(postgres.String(id.String())).AS_UUID()
	}

	var dest []struct {
		CardId uuid.UUID `alias:"account_card.card_id"`
		Level  int32     `alias:"account_card.level"`
	}

	stmt := postgres.SELECT(AccountCard.CardID, AccountCard.Level).
		FROM(AccountCard).
		WHERE(
			AccountCard.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()).
				AND(AccountCard.CardID.IN(expressions...)),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	result := make(map[uuid.UUID]OwnedCardInfo, len(dest))
	for _, d := range dest {
		result[d.CardId] = OwnedCardInfo{CardId: d.CardId, Level: d.Level}
	}

	return result, nil
}
