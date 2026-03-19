package account

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type UpdateOneParam struct {
	Username string
	Email    string
	Avatar   *string
}

func (r *Repository) UpdateOneByDiscordId(ctx context.Context, discordId string, param UpdateOneParam) (model.Account, error) {
	var dest model.Account

	stmt := Account.UPDATE(
		Account.Username,
		Account.Email,
		Account.Avatar,
	).SET(
		param.Username,
		param.Email,
		param.Avatar,
	).WHERE(
		Account.DiscordId.EQ(postgres.String(discordId)),
	).RETURNING(Account.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
