package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

// FindByPool returns all packs in a pool, ordered by last_pool_activated_at ASC NULLS FIRST.
// This ensures never-activated packs are selected first for rotation.
func (r *Repository) FindByPool(ctx context.Context, poolId string) ([]model.Pack, error) {
	var dest []model.Pack

	stmt := postgres.SELECT(Pack.AllColumns).
		FROM(Pack).
		WHERE(Pack.PoolID.EQ(postgres.CAST(postgres.String(poolId)).AS_UUID())).
		ORDER_BY(Pack.LastPoolActivatedAt.ASC().NULLS_FIRST())

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
