package profile

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type UpsertParam struct {
	AccountId     string
	DisplayName   string
	Bio           *string
	ShowcaseCards string
}

func (r *Repository) Upsert(ctx context.Context, param UpsertParam) (model.Profile, error) {
	var dest model.Profile

	stmt := Profile.INSERT(
		Profile.AccountID,
		Profile.DisplayName,
		Profile.Bio,
		Profile.ShowcaseCards,
	).VALUES(
		param.AccountId,
		param.DisplayName,
		param.Bio,
		param.ShowcaseCards,
	).ON_CONFLICT(Profile.AccountID).DO_UPDATE(
		postgres.SET(
			Profile.DisplayName.SET(Profile.EXCLUDED.DisplayName),
			Profile.Bio.SET(Profile.EXCLUDED.Bio),
			Profile.ShowcaseCards.SET(Profile.EXCLUDED.ShowcaseCards),
			Profile.UpdatedAt.SET(postgres.NOW()),
		),
	).RETURNING(Profile.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return dest, trace.Wrap(err)
	}

	return dest, nil
}
