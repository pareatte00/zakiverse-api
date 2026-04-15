package profile

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOneByAccountId(ctx context.Context, accountId string) (model.Profile, error) {
	var dest model.Profile

	stmt := postgres.SELECT(Profile.AllColumns).
		FROM(Profile).
		WHERE(Profile.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}

func (r *Repository) FindOneByDisplayName(ctx context.Context, displayName string) (model.Profile, error) {
	var dest model.Profile

	stmt := postgres.SELECT(Profile.AllColumns).
		FROM(Profile).
		WHERE(postgres.LOWER(Profile.DisplayName).EQ(postgres.LOWER(postgres.String(displayName))))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
