package account_card

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	AccountId string
	CardId    string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.AccountCard, error) {
	var dest model.AccountCard

	stmt := AccountCard.INSERT(
		AccountCard.AccountID,
		AccountCard.CardID,
	).VALUES(
		param.AccountId,
		param.CardId,
	).RETURNING(AccountCard.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
