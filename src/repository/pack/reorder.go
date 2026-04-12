package pack

import (
	"context"
	"fmt"

	"github.com/zakiverse/zakiverse-api/util/trace"
)

// Reorder updates sort_order for each pack based on its index in the ids slice.
func (r *Repository) Reorder(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	query := "UPDATE pack SET sort_order = CASE id"
	args := make([]any, 0, len(ids)*2+len(ids))

	for i, id := range ids {
		query += fmt.Sprintf(" WHEN $%d::uuid THEN $%d", i*2+1, i*2+2)
		args = append(args, id, i)
	}
	query += " END WHERE id IN ("
	for i, id := range ids {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("$%d::uuid", len(ids)*2+i+1)
		args = append(args, id)
	}
	query += ")"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
