package anime

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

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]model.Anime, error) {
	var dest []model.Anime

	condition := postgres.Bool(true)
	if param.Search != "" {
		search := postgres.String("%" + param.Search + "%")
		condition = condition.AND(postgres.LOWER(Anime.Title).LIKE(postgres.LOWER(search)))
	}

	stmt := postgres.SELECT(Anime.AllColumns).
		FROM(Anime).
		WHERE(condition).
		ORDER_BY(Anime.Title.ASC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
