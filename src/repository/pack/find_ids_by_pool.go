package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindIdsByPoolId(ctx context.Context, poolId string) ([]uuid.UUID, error) {
	var dest []struct {
		ID uuid.UUID
	}

	stmt := postgres.SELECT(Pack.ID).
		FROM(Pack).
		WHERE(Pack.PoolID.EQ(postgres.CAST(postgres.String(poolId)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	ids := make([]uuid.UUID, len(dest))
	for i, d := range dest {
		ids[i] = d.ID
	}

	return ids, nil
}
