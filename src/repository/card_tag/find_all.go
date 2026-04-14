package card_tag

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllParam struct {
	Search string
	Limit  int64
	Offset int64
}

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]model.CardTag, error) {
	var dest []model.CardTag

	var where postgres.BoolExpression
	if param.Search != "" {
		where = CardTag.Name.LIKE(postgres.String("%" + param.Search + "%"))
	}

	stmt := postgres.SELECT(CardTag.AllColumns).
		FROM(CardTag)

	if where != nil {
		stmt = stmt.WHERE(where)
	}

	stmt = stmt.ORDER_BY(CardTag.Name.ASC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
