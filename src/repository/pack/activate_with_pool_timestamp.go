package pack

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

// ActivateWithPoolTimestamp sets is_active=true and last_pool_activated_at=now for a pack.
func (r *Repository) ActivateWithPoolTimestamp(ctx context.Context, id string, now time.Time) error {
	stmt := Pack.UPDATE(Pack.IsActive, Pack.LastPoolActivatedAt).
		SET(true, now).
		WHERE(Pack.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()))

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
