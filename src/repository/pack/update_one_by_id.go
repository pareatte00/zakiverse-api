package pack

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) UpdateOneById(ctx context.Context, id string, updates map[string]any) (model.Pack, error) {
	var dest model.Pack

	columnMap := map[string]postgres.Column{
		"code":           Pack.Code,
		"name":           Pack.Name,
		"description":    Pack.Description,
		"image":          Pack.Image,
		"name_image":     Pack.NameImage,
		"cards_per_pull": Pack.CardsPerPull,
		"sort_order":     Pack.SortOrder,
		"config":         Pack.Config,
		"rotation_order": Pack.RotationOrder,
	}

	var f postgres.ColumnList
	var vs []any
	for k, v := range updates {
		f = append(f, columnMap[k])
		vs = append(vs, v)
	}

	stmt := Pack.UPDATE(f).SET(vs[0], vs[1:]...).WHERE(
		Pack.ID.EQ(postgres.CAST(postgres.String(id)).AS_UUID()),
	).RETURNING(Pack.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
