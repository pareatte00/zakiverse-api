package pack

import (
	"context"
	"time"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	Code         string
	Name         string
	Description  *string
	Image        string
	NameImage    *string
	Type         string
	CardsPerPull int32
	SortOrder    int32
	IsActive     bool
	OpenAt       *time.Time
	CloseAt      *time.Time
	Config       string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.Pack, error) {
	var dest model.Pack

	stmt := Pack.INSERT(
		Pack.Code,
		Pack.Name,
		Pack.Description,
		Pack.Image,
		Pack.NameImage,
		Pack.Type,
		Pack.CardsPerPull,
		Pack.SortOrder,
		Pack.IsActive,
		Pack.OpenAt,
		Pack.CloseAt,
		Pack.Config,
	).VALUES(
		param.Code,
		param.Name,
		param.Description,
		param.Image,
		param.NameImage,
		param.Type,
		param.CardsPerPull,
		param.SortOrder,
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
