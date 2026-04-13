package pack

import (
	"context"
	"fmt"
	"strings"

	"github.com/zakiverse/zakiverse-api/util/trace"
)

// AssignToPool sets pool_id for the given pack IDs and NULLs pool_id for any packs
// currently in the pool but not in the list. An empty ids list unassigns all packs.
func (r *Repository) AssignToPool(ctx context.Context, poolId string, ids []string) error {
	if len(ids) == 0 {
		// Unassign all packs from this pool
		query := "UPDATE pack SET pool_id = NULL WHERE pool_id = $1::uuid"
		_, err := r.db.ExecContext(ctx, query, poolId)
		if err != nil {
			return trace.Wrap(err)
		}
		return nil
	}

	// Build placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]any, 0, len(ids)+1)
	args = append(args, poolId) // $1 = poolId
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d::uuid", i+2)
		args = append(args, id)
	}
	inClause := strings.Join(placeholders, ", ")

	// Assign: set pool_id for listed packs
	assignQuery := fmt.Sprintf(
		"UPDATE pack SET pool_id = $1::uuid WHERE id IN (%s)",
		inClause,
	)
	_, err := r.db.ExecContext(ctx, assignQuery, args...)
	if err != nil {
		return trace.Wrap(err)
	}

	// Unassign: NULL pool_id for packs in this pool but not in the list
	unassignQuery := fmt.Sprintf(
		"UPDATE pack SET pool_id = NULL WHERE pool_id = $1::uuid AND id NOT IN (%s)",
		inClause,
	)
	_, err = r.db.ExecContext(ctx, unassignQuery, args...)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
