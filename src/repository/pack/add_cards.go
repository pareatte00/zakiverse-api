package pack

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type AddCardParam struct {
	PackId       string
	CardId       string
	Weight       float64
	IsFeatured   bool
	FeaturedRate *float64
}

func (r *Repository) AddCards(ctx context.Context, params []AddCardParam) ([]model.PackCard, error) {
	var dest []model.PackCard

	stmt := PackCard.INSERT(
		PackCard.PackID,
		PackCard.CardID,
		PackCard.Weight,
		PackCard.IsFeatured,
		PackCard.FeaturedRate,
	)

	for _, p := range params {
		stmt = stmt.VALUES(p.PackId, p.CardId, p.Weight, p.IsFeatured, p.FeaturedRate)
	}

	err := stmt.RETURNING(PackCard.AllColumns).QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
