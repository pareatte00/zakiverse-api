package pack_pool

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) UpdateOneById(ctx context.Context, id string, updates map[string]any) (model.PackPool, error) {
	var dest model.PackPool

	columnMap := map[string]postgres.Column{
		"name":                PackPool.Name,
		"description":         PackPool.Description,
		"image":               PackPool.Image,
		"banner_type":         PackPool.BannerType,
		"sort_order":          PackPool.SortOrder,
		"is_active":           PackPool.IsActive,
		"open_at":             PackPool.OpenAt,
		"close_at":            PackPool.CloseAt,
		"active_count":        PackPool.ActiveCount,
		"rotation_type":       PackPool.RotationType,
		"rotation_day":        PackPool.RotationDay,
		"rotation_interval":   PackPool.RotationInterval,
		"rotation_hour":       PackPool.RotationHour,
		"rotation_order_mode": PackPool.RotationOrderMode,
		"next_rotation_at":    PackPool.NextRotationAt,
		"preview_days":        PackPool.PreviewDays,
	}

	var f postgres.ColumnList
	var vs []any
	for k, v := range updates {
		f = append(f, columnMap[k])
		vs = append(vs, v)
	}

	stmt := PackPool.UPDATE(f).SET(vs[0], vs[1:]...).WHERE(
		PackPool.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()),
	).RETURNING(PackPool.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
