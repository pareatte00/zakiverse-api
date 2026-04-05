package pack

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindPendingOpen(ctx context.Context, now time.Time) ([]model.Pack, error) {
	var dest []model.Pack

	stmt := postgres.SELECT(Pack.AllColumns).
		FROM(Pack).
		WHERE(
			Pack.IsActive.EQ(postgres.Bool(false)).
				AND(Pack.PoolID.IS_NULL()).
				AND(Pack.OpenAt.IS_NOT_NULL()).
				AND(Pack.OpenAt.LT_EQ(postgres.TimestampzT(now))),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}

func (r *Repository) FindExpired(ctx context.Context, now time.Time) ([]model.Pack, error) {
	var dest []model.Pack

	stmt := postgres.SELECT(Pack.AllColumns).
		FROM(Pack).
		WHERE(
			Pack.IsActive.EQ(postgres.Bool(true)).
				AND(Pack.PoolID.IS_NULL()).
				AND(Pack.CloseAt.IS_NOT_NULL()).
				AND(Pack.CloseAt.LT_EQ(postgres.TimestampzT(now))),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}

func (r *Repository) SetActive(ctx context.Context, id string, active bool) error {
	stmt := Pack.UPDATE(Pack.IsActive).SET(active).WHERE(
		Pack.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()),
	)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
