package account_pull_history

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PullHistoryWithCard struct {
	model.AccountPullHistory

	CardName  string  `alias:"card.name"`
	CardImage string  `alias:"card.image"`
	TagName   *string `alias:"card_tag.name"`
}

type FindByAccountAndPackParam struct {
	AccountId string
	PackId    string
	Limit     int64
	Offset    int64
}

func (r *Repository) FindByAccountAndPack(ctx context.Context, param FindByAccountAndPackParam) ([]PullHistoryWithCard, error) {
	var dest []PullHistoryWithCard

	stmt := postgres.SELECT(
		AccountPullHistory.AllColumns,
		Card.Name,
		Card.Image,
		CardTag.Name,
	).FROM(
		AccountPullHistory.
			INNER_JOIN(Card, Card.ID.EQ(AccountPullHistory.CardID)).
			LEFT_JOIN(CardTag, CardTag.ID.EQ(Card.TagID)),
	).WHERE(
		AccountPullHistory.AccountID.EQ(postgres.CAST(postgres.String(param.AccountId)).AS_UUID()).
			AND(AccountPullHistory.PackID.EQ(postgres.CAST(postgres.String(param.PackId)).AS_UUID())),
	).ORDER_BY(
		AccountPullHistory.PulledAt.DESC(),
	).LIMIT(param.Limit).
		OFFSET(param.Offset)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
