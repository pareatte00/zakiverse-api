package card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllByAnimeIdParam struct {
	AnimeId string
	Limit   int64
	Offset  int64
}

func (r *Repository) FindAllByAnimeId(ctx context.Context, param FindAllByAnimeIdParam) ([]CardWithAnime, error) {
	var dest []CardWithAnime

	stmt := postgres.SELECT(Card.AllColumns, Anime.AllColumns).
		FROM(Card.INNER_JOIN(Anime, Anime.ID.EQ(Card.AnimeID))).
		WHERE(Card.AnimeID.EQ(postgres.CAST(postgres.String(param.AnimeId)).AS_UUID())).
		ORDER_BY(Card.Name.ASC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
