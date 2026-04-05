package pack_pool

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

// FindDueRotations returns pools where rotation_day matches today's weekday
// and last_rotated_at is NULL or before the start of the current week.
func (r *Repository) FindDueRotations(ctx context.Context, now time.Time) ([]model.PackPool, error) {
	var dest []model.PackPool

	weekday := int32(now.Weekday()) // 0=Sun..6=Sat

	// Start of current week (Sunday 00:00)
	daysFromSunday := int(now.Weekday())
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-daysFromSunday, 0, 0, 0, 0, now.Location())

	stmt := postgres.SELECT(PackPool.AllColumns).
		FROM(PackPool).
		WHERE(
			PackPool.RotationDay.EQ(postgres.Int32(weekday)).AND(
				PackPool.LastRotatedAt.IS_NULL().OR(
					PackPool.LastRotatedAt.LT(postgres.TimestampzT(weekStart)),
				),
			),
		)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
