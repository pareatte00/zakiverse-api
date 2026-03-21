package account

import (
	"context"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CreateOneParam struct {
	DiscordId string
	Username  string
	Email     string
	Avatar    *string
}

func (r *Repository) CreateOne(ctx context.Context, param CreateOneParam) (model.Account, error) {
	var dest model.Account

	stmt := Account.INSERT(
		Account.DiscordID,
		Account.Username,
		Account.Email,
		Account.Avatar,
	).VALUES(
		param.DiscordId,
		param.Username,
		param.Email,
		param.Avatar,
	).RETURNING(Account.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
