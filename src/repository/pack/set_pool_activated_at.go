package pack

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) SetPoolActivatedAt(ctx context.Context, id string, now time.Time) error {
	stmt := Pack.UPDATE(Pack.LastPoolActivatedAt).
		SET(now).
		WHERE(Pack.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()))

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
