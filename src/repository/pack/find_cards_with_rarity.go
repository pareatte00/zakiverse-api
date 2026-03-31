package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PackCardWithRarity struct {
	model.PackCard

	Card model.Card
}

func (r *Repository) FindCardsWithRarity(ctx context.Context, packId string) ([]PackCardWithRarity, error) {
	var dest []PackCardWithRarity

	stmt := postgres.SELECT(PackCard.AllColumns, Card.Rarity).
		FROM(PackCard.INNER_JOIN(Card, Card.ID.EQ(PackCard.CardID))).
		WHERE(PackCard.PackID.EQ(postgres.CAST(postgres.String(packId)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
