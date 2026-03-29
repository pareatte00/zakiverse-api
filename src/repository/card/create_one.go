package card

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	MalId   int32
	AnimeId string
	Rarity  string
	Name    string
	Image   string
	Config  string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.Card, error) {
	var dest model.Card

	stmt := Card.INSERT(
		Card.MalID,
		Card.AnimeID,
		Card.Rarity,
		Card.Name,
		Card.Image,
		Card.Config,
	).VALUES(
		param.MalId,
		param.AnimeId,
		param.Rarity,
		param.Name,
		param.Image,
		param.Config,
	).RETURNING(Card.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
