package anime

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var dest struct {
		Count int64
	}

	stmt := postgres.SELECT(
		postgres.COUNT(postgres.STAR).AS("count"),
	).FROM(Anime)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}
