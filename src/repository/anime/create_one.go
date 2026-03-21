package anime

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	MalId      int32
	Title      string
	Synopsis   *string
	CoverImage *string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.Anime, error) {
	var dest model.Anime

	stmt := Anime.INSERT(
		Anime.MalID,
		Anime.Title,
		Anime.Synopsis,
		Anime.CoverImage,
	).VALUES(
		param.MalId,
		param.Title,
		param.Synopsis,
		param.CoverImage,
	).RETURNING(Anime.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
