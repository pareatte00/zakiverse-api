package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

// FindCurrentByPool returns the currently active packs in a pool.
// For rotating pools, returns the most recently activated `activeCount` packs.
// For non-rotating pools (activeCount <= 0), returns all packs.
func (r *Repository) FindCurrentByPool(ctx context.Context, poolId string, activeCount int32) ([]PackWithCardCount, error) {
	var dest []PackWithCardCount

	stmt := postgres.SELECT(
		Pack.AllColumns,
		postgres.COUNT(PackCard.ID).AS("pack_with_card_count.total_cards"),
	).
		FROM(Pack.LEFT_JOIN(PackCard, PackCard.PackID.EQ(Pack.ID))).
		WHERE(Pack.PoolID.EQ(postgres.CAST(postgres.String(poolId)).AS_UUID())).
		GROUP_BY(Pack.ID).
		ORDER_BY(Pack.LastPoolActivatedAt.DESC().NULLS_LAST())

	if activeCount > 0 {
		stmt = stmt.LIMIT(int64(activeCount))
	}

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
