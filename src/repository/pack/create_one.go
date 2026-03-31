package pack

import (
	"context"
	"time"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Name         string
	Description  *string
	Image        string
	CardsPerPull int32
	IsActive     bool
	OpenAt       *time.Time
	CloseAt      *time.Time
	Config       string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.Pack, error) {
	var dest model.Pack

	stmt := Pack.INSERT(
		Pack.Name,
		Pack.Description,
		Pack.Image,
		Pack.CardsPerPull,
		Pack.IsActive,
		Pack.OpenAt,
		Pack.CloseAt,
		Pack.Config,
	).VALUES(
		param.Name,
		param.Description,
		param.Image,
		param.CardsPerPull,
		param.IsActive,
		param.OpenAt,
		param.CloseAt,
		param.Config,
	).RETURNING(Pack.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
