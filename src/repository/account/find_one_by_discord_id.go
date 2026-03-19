package account

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

func (r *Repository) FindOneByDiscordId(ctx context.Context, discordId string) (model.Account, error) {
	var dest model.Account

	stmt := postgres.SELECT(Account.AllColumns).
		FROM(Account).
		WHERE(Account.DiscordId.EQ(postgres.String(discordId)))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
