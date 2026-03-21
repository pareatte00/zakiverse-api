package anime

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type UpdateOneByIdParam struct {
	Title      string
	Synopsis   *string
	CoverImage *string
}

func (r *Repository) UpdateOneById(ctx context.Context, id string, param UpdateOneByIdParam) (model.Anime, error) {
	var dest model.Anime

	stmt := Anime.UPDATE(
		Anime.Title,
		Anime.Synopsis,
		Anime.CoverImage,
	).SET(
		param.Title,
		param.Synopsis,
		param.CoverImage,
	).WHERE(
		Anime.ID.EQ(postgres.String(id)),
	).RETURNING(Anime.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
