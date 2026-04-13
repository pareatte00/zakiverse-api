package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CountParam struct {
	Search     string
	Unassigned bool
}

func (r *Repository) Count(ctx context.Context, param CountParam) (int64, error) {
	var dest struct {
		Count int64
	}

	condition := postgres.Bool(true)
	if param.Search != "" {
		search := postgres.String("%" + param.Search + "%")
		condition = condition.AND(postgres.LOWER(Pack.Name).LIKE(postgres.LOWER(search)))
	}
	if param.Unassigned {
		condition = condition.AND(Pack.PoolID.IS_NULL())
	}

	stmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("count"),
	).FROM(Pack).WHERE(condition)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}
