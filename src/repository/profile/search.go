package profile

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type SearchParam struct {
	Query  string
	Limit  int64
	Offset int64
}

func (r *Repository) Search(ctx context.Context, param SearchParam) ([]model.Profile, error) {
	var dest []model.Profile

	search := postgres.String("%" + param.Query + "%")
	stmt := postgres.SELECT(Profile.AllColumns).
		FROM(Profile).
		WHERE(postgres.LOWER(Profile.DisplayName).LIKE(postgres.LOWER(search))).
		ORDER_BY(Profile.DisplayName.ASC()).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}

func (r *Repository) CountSearch(ctx context.Context, query string) (int64, error) {
	var dest struct {
		Count int64
	}

	search := postgres.String("%" + query + "%")
	stmt := postgres.SELECT(postgres.COUNT(Profile.AccountID).AS("count")).
		FROM(Profile).
		WHERE(postgres.LOWER(Profile.DisplayName).LIKE(postgres.LOWER(search)))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}
