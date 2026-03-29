package anime

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) UpdateOneById(ctx context.Context, id string, updates map[string]any) (model.Anime, error) {
	var dest model.Anime

	columnMap := map[string]postgres.Column{
		"title":       Anime.Title,
		"synopsis":    Anime.Synopsis,
		"cover_image": Anime.CoverImage,
	}

	var f postgres.ColumnList
	var vs []any
	for k, v := range updates {
		f = append(f, columnMap[k])
		vs = append(vs, v)
	}

	stmt := Anime.UPDATE(f).SET(vs[0], vs[1:]...).WHERE(
		Anime.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()),
	).RETURNING(Anime.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
