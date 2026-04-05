package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

// DeactivateByPool sets is_active=false for all active packs in a pool.
func (r *Repository) DeactivateByPool(ctx context.Context, poolId string) error {
	stmt := Pack.UPDATE(Pack.IsActive).
		SET(false).
		WHERE(
			Pack.PoolID.EQ(postgres.CAST(postgres.String(poolId)).AS_UUID()).
				AND(Pack.IsActive.EQ(postgres.Bool(true))),
		)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
