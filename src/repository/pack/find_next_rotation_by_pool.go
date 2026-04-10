package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

// FindNextRotationByPool returns the packs that would be activated next in a rotation.
// Uses the same ordering as rotation (ASC NULLS FIRST) and takes the first activeCount.
func (r *Repository) FindNextRotationByPool(ctx context.Context, poolId string, activeCount int32) ([]model.Pack, error) {
	var dest []model.Pack

	stmt := postgres.SELECT(Pack.AllColumns).
		FROM(Pack).
		WHERE(Pack.PoolID.EQ(postgres.CAST(postgres.String(poolId)).AS_UUID())).
		ORDER_BY(Pack.LastPoolActivatedAt.ASC().NULLS_FIRST()).
		LIMIT(int64(activeCount))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
