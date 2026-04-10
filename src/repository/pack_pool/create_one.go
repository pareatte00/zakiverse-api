package pack_pool

import (
	"context"
	"time"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Name              string
	Description       *string
	ActiveCount       int32
	RotationDay       *int32
	Image             *string
	BannerType        string
	SortOrder         int32
	IsActive          bool
	OpenAt            *time.Time
	CloseAt           *time.Time
	RotationType      string
	RotationInterval  int32
	RotationHour      int32
	RotationOrderMode string
	NextRotationAt    *time.Time
	PreviewDays       int32
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.PackPool, error) {
	var dest model.PackPool

	stmt := PackPool.INSERT(
		PackPool.Name,
		PackPool.Description,
		PackPool.ActiveCount,
		PackPool.RotationDay,
		PackPool.Image,
		PackPool.BannerType,
		PackPool.SortOrder,
		PackPool.IsActive,
		PackPool.OpenAt,
		PackPool.CloseAt,
		PackPool.RotationType,
		PackPool.RotationInterval,
		PackPool.RotationHour,
		PackPool.RotationOrderMode,
		PackPool.NextRotationAt,
		PackPool.PreviewDays,
	).VALUES(
		param.Name,
		param.Description,
		param.ActiveCount,
		param.RotationDay,
		param.Image,
		param.BannerType,
		param.SortOrder,
		param.IsActive,
		param.OpenAt,
		param.CloseAt,
		param.RotationType,
		param.RotationInterval,
		param.RotationHour,
		param.RotationOrderMode,
		param.NextRotationAt,
		param.PreviewDays,
	).RETURNING(PackPool.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
