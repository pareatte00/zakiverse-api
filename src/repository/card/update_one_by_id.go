package card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type UpdateOneByIdParam struct {
	RarityId string
	Name     string
	Image    string
	Config   string
}

func (r *Repository) UpdateOneById(ctx context.Context, id string, param UpdateOneByIdParam) (model.Card, error) {
	var dest model.Card

	stmt := Card.UPDATE(
		Card.RarityID,
		Card.Name,
		Card.Image,
		Card.Config,
	).SET(
		param.RarityId,
		param.Name,
		param.Image,
		param.Config,
	).WHERE(
		Card.ID.EQ(postgres.String(id)),
	).RETURNING(Card.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
