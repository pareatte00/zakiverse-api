package card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOneById(ctx context.Context, id string) (model.Card, error) {
	var dest model.Card

	stmt := postgres.SELECT(Card.AllColumns).
		FROM(Card).
		WHERE(Card.ID.EQ(postgres.String(id)))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
