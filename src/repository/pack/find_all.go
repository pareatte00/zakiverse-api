package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllParam struct {
	Search     string
	Unassigned bool
	Limit      int64
	Offset     int64
}

type PackWithCardCount struct {
	model.Pack

	TotalCards int64
}

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]PackWithCardCount, error) {
	var dest []PackWithCardCount

	condition := postgres.Bool(true)
	if param.Search != "" {
		search := postgres.String("%" + param.Search + "%")
		condition = condition.AND(postgres.LOWER(Pack.Name).LIKE(postgres.LOWER(search)))
	}
	if param.Unassigned {
		condition = condition.AND(Pack.PoolID.IS_NULL())
	}

	stmt := postgres.SELECT(
		Pack.AllColumns,
		postgres.COUNT(PackCard.ID).AS("pack_with_card_count.total_cards"),
	).
		FROM(Pack.LEFT_JOIN(PackCard, PackCard.PackID.EQ(Pack.ID))).
		WHERE(condition).
		GROUP_BY(Pack.ID).
		ORDER_BY(Pack.SortOrder.ASC(), Pack.CreatedAt.DESC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
