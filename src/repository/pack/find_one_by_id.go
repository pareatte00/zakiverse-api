package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOneById(ctx context.Context, id string) (PackWithCards, error) {
	var dest PackWithCards

	stmt := postgres.SELECT(Pack.AllColumns, PackCard.AllColumns, Card.AllColumns, Anime.AllColumns).
		FROM(
			Pack.
				LEFT_JOIN(PackCard, PackCard.PackID.EQ(Pack.ID)).
				LEFT_JOIN(Card, Card.ID.EQ(PackCard.CardID)).
				LEFT_JOIN(Anime, Anime.ID.EQ(Card.AnimeID)),
		).
		WHERE(Pack.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
