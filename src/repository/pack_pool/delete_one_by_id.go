package pack_pool

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) DeleteOneById(ctx context.Context, id string) error {
	stmt := PackPool.DELETE().WHERE(
		PackPool.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()),
	)

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
