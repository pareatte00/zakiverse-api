package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) RemoveCards(ctx context.Context, packId string, cardIds []string) error {
	ids := make([]postgres.Expression, len(cardIds))
	for i, id := range cardIds {
		ids[i] = postgres.CAST(postgres.String(id)).AS_UUID()
	}

	stmt := PackCard.DELETE().WHERE(
		PackCard.PackID.EQ(postgres.CAST(postgres.String(packId)).AS_UUID()).
			AND(PackCard.CardID.IN(ids...)),
	)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
