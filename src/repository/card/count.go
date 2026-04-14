package card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CountParam struct {
	Search  string
	Rarity  string
	TagId   string
	AnimeId string
}

func (r *Repository) Count(ctx context.Context, param CountParam) (int64, error) {
	var dest struct {
		Count int64
	}

	condition := postgres.Bool(true)

	if param.Search != "" {
		search := postgres.String("%" + param.Search + "%")
		condition = condition.AND(
			postgres.LOWER(Card.Name).LIKE(postgres.LOWER(search)).
				OR(postgres.LOWER(Anime.Title).LIKE(postgres.LOWER(search))),
		)
	}

	if param.Rarity != "" {
		condition = condition.AND(Card.Rarity.EQ(postgres.NewEnumValue(param.Rarity)))
	}

	if param.TagId != "" {
		condition = condition.AND(Card.TagID.EQ(postgres.CAST(postgres.String(param.TagId)).AS_UUID()))
	}

	if param.AnimeId != "" {
		condition = condition.AND(Card.AnimeID.EQ(postgres.CAST(postgres.String(param.AnimeId)).AS_UUID()))
	}

	stmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("count"),
	).FROM(
		Card.INNER_JOIN(Anime, Anime.ID.EQ(Card.AnimeID)),
	).WHERE(condition)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}

