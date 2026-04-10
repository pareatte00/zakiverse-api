package card_tag

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Name string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.CardTag, error) {
	var dest model.CardTag

	stmt := CardTag.INSERT(
		CardTag.Name,
	).VALUES(
		param.Name,
	).RETURNING(CardTag.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
