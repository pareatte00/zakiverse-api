package profile

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/table"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CompletionStat struct {
	Rarity string
	Owned  int64
	Total  int64
}

func (r *Repository) GetCompletionStats(ctx context.Context, accountId string) ([]CompletionStat, error) {
	var dest []CompletionStat

	stmt := postgres.RawStatement(`
		SELECT
			c.rarity AS "rarity",
			COUNT(DISTINCT ac.card_id) AS "owned",
			COUNT(DISTINCT c.id) AS "total"
		FROM card c
		LEFT JOIN account_card ac ON ac.card_id = c.id AND ac.account_id = :account_id::uuid
		GROUP BY c.rarity
	`, postgres.RawArgs{":account_id": accountId})

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return dest, nil
}

type HighestLevelCard struct {
	CardId uuid.UUID
	Name   string
	Image  string
	Rarity string
	Level  int32
}

func (r *Repository) GetHighestLevelCard(ctx context.Context, accountId string) (*HighestLevelCard, error) {
	var dest HighestLevelCard

	stmt := postgres.RawStatement(`
		SELECT
			c.id AS "card_id",
			c.name AS "name",
			c.image AS "image",
			c.rarity AS "rarity",
			ac.level AS "level"
		FROM account_card ac
		JOIN card c ON c.id = ac.card_id
		WHERE ac.account_id = :account_id::uuid
		ORDER BY ac.level DESC, ac.obtained_at ASC
		LIMIT 1
	`, postgres.RawArgs{":account_id": accountId})

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &dest, nil
}

func (r *Repository) CountCards(ctx context.Context, accountId string) (int64, error) {
	var dest struct {
		Count int64
	}

	stmt := postgres.SELECT(postgres.COUNT(AccountCard.ID).AS("count")).
		FROM(AccountCard).
		WHERE(AccountCard.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}

func (r *Repository) CountPulls(ctx context.Context, accountId string) (int64, error) {
	var dest struct {
		Count int64
	}

	stmt := postgres.SELECT(postgres.COUNT(AccountPullHistory.ID).AS("count")).
		FROM(AccountPullHistory).
		WHERE(AccountPullHistory.AccountID.EQ(postgres.CAST(postgres.String(accountId)).AS_UUID()))

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return 0, trace.Wrap(err)
	}

	return dest.Count, nil
}

type RecentPull struct {
	CardId   uuid.UUID
	Name     string
	Image    string
	Rarity   string
	PulledAt string
	IsNew    bool
}

func (r *Repository) GetRecentPulls(ctx context.Context, accountId string, limit int) ([]RecentPull, error) {
	var dest []RecentPull

	stmt := postgres.RawStatement(`
		SELECT
			c.id AS "card_id",
			c.name AS "name",
			c.image AS "image",
			aph.rarity AS "rarity",
			aph.pulled_at AS "pulled_at",
			aph.is_new AS "is_new"
		FROM account_pull_history aph
		JOIN card c ON c.id = aph.card_id
		WHERE aph.account_id = :account_id::uuid
		ORDER BY aph.pulled_at DESC
		LIMIT :limit
	`, postgres.RawArgs{":account_id": accountId, ":limit": limit})

	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return dest, nil
}
