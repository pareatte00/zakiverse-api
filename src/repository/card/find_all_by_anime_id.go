package card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllByAnimeIdParam struct {
	AnimeId string
	Limit   int64
	Offset  int64
}

func (r *Repository) FindAllByAnimeId(ctx context.Context, param FindAllByAnimeIdParam) ([]model.Card, error) {
	var dest []model.Card

	stmt := postgres.SELECT(Card.AllColumns).
		FROM(Card).
		WHERE(Card.AnimeID.EQ(postgres.String(param.AnimeId))).
		ORDER_BY(Card.Name.ASC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
