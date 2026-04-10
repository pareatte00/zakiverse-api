package card

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type FindAllParam struct {
	Search string
	Rarity string
	TagId  string
	Sort   string
	Order  string
	Limit  int64
	Offset int64
}

func (r *Repository) FindAll(ctx context.Context, param FindAllParam) ([]CardWithAnime, error) {
	var dest []CardWithAnime

	condition := postgres.Bool(true)

	if param.Search != "" {
		search := postgres.String("%" + param.Search + "%")
		condition = condition.AND(
			postgres.LOWER(Card.Name).LIKE(postgres.LOWER(search)).
				OR(postgres.LOWER(Anime.Title).LIKE(postgres.LOWER(search))),
		)
	}

	if param.Rarity != "" {
		condition = condition.AND(Card.Rarity.EQ(postgres.NewEnumValue(param.Rarity)))
	}

	if param.TagId != "" {
		condition = condition.AND(Card.TagID.EQ(postgres.CAST(postgres.String(param.TagId)).AS_UUID()))
	}

	sortColumn := map[string]postgres.Column{
		"name":   Card.Name,
		"rarity": Card.Rarity,
	}

	col, ok := sortColumn[param.Sort]
	if !ok {
		col = Card.Name
	}

	var orderClauses []postgres.OrderByClause
	if param.Order == "desc" {
		orderClauses = append(orderClauses, col.DESC())
	} else {
		orderClauses = append(orderClauses, col.ASC())
	}
	if param.Sort == "rarity" {
		orderClauses = append(orderClauses, Card.Name.ASC())
	}

	stmt := postgres.SELECT(Card.AllColumns, Anime.AllColumns, CardTag.AllColumns).
		FROM(
			Card.INNER_JOIN(Anime, Anime.ID.EQ(Card.AnimeID)).
				LEFT_JOIN(CardTag, CardTag.ID.EQ(Card.TagID)),
		).
		WHERE(condition).
		ORDER_BY(orderClauses...).
		LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
