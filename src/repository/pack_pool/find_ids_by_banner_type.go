package pack_pool

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindIdsByBannerType(ctx context.Context, bannerType string) ([]uuid.UUID, error) {
	var dest []struct {
		ID uuid.UUID
	}

	stmt := postgres.SELECT(PackPool.ID).
		FROM(PackPool).
		WHERE(PackPool.BannerType.EQ(postgres.NewEnumValue(bannerType)))

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
