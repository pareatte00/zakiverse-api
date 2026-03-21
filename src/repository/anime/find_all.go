package anime

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllParam struct {
	Limit  int64
	Offset int64
}

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]model.Anime, error) {
	var dest []model.Anime

	stmt := postgres.SELECT(Anime.AllColumns).
		FROM(Anime).
		ORDER_BY(Anime.Title.ASC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
