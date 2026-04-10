package card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOneById(ctx context.Context, id string) (CardWithAnime, error) {
	var dest CardWithAnime

	stmt := postgres.SELECT(Card.AllColumns, Anime.AllColumns, CardTag.AllColumns).
		FROM(
			Card.INNER_JOIN(Anime, Anime.ID.EQ(Card.AnimeID)).
				LEFT_JOIN(CardTag, CardTag.ID.EQ(Card.TagID)),
		).
		WHERE(Card.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
