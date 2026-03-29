package account

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) UpdateOneByDiscordId(ctx context.Context, discordId string, updates map[string]any) (model.Account, error) {
	var dest model.Account

	columnMap := map[string]postgres.Column{
		"username": Account.Username,
		"email":    Account.Email,
		"avatar":   Account.Avatar,
	}

	var f postgres.ColumnList
	var vs []any
	for k, v := range updates {
		f = append(f, columnMap[k])
		vs = append(vs, v)
	}

	stmt := Account.UPDATE(f).SET(vs[0], vs[1:]...).WHERE(
		Account.DiscordID.EQ(postgres.String(discordId)),
	).RETURNING(Account.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
