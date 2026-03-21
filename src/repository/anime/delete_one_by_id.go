package anime

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) DeleteOneById(ctx context.Context, id string) error {
	stmt := Anime.DELETE().
		WHERE(Anime.ID.EQ(postgres.String(id)))

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
