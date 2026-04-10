package repository

import (
	"database/sql"

	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/repository/account"
	"github.com/zakiverse/zakiverse-api/src/repository/account_card"
	"github.com/zakiverse/zakiverse-api/src/repository/account_pack_pity"
	"github.com/zakiverse/zakiverse-api/src/repository/account_pull_history"
	"github.com/zakiverse/zakiverse-api/src/repository/anime"
	"github.com/zakiverse/zakiverse-api/src/repository/card"
	"github.com/zakiverse/zakiverse-api/src/repository/card_tag"
	"github.com/zakiverse/zakiverse-api/src/repository/pack"
	"github.com/zakiverse/zakiverse-api/src/repository/pack_pool"
)

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Database   *sql.DB
}

type Repository struct {
	config     config.ConfigConstant
	credential config.ConfigCredential

	Account            *account.Repository
	Anime              *anime.Repository
	Card               *card.Repository
	CardTag            *card_tag.Repository
	AccountCard        *account_card.Repository
	Pack               *pack.Repository
	AccountPackPity    *account_pack_pity.Repository
	AccountPullHistory *account_pull_history.Repository
	PackPool           *pack_pool.Repository
}

func New(d Dependency) *Repository {
	return &Repository{
		config:     d.Config,
		credential: d.Credential,

		Account:            account.New(d.Database),
		Anime:              anime.New(d.Database),
		Card:               card.New(d.Database),
		CardTag:            card_tag.New(d.Database),
		AccountCard:        account_card.New(d.Database),
		Pack:               pack.New(d.Database),
		AccountPackPity:    account_pack_pity.New(d.Database),
		AccountPullHistory: account_pull_history.New(d.Database),
		PackPool:           pack_pool.New(d.Database),
	}
}

func (r *Repository) Tx(tx *sql.Tx) *Repository {
	return &Repository{
		config:     r.config,
		credential: r.credential,

		Account:            account.New(tx),
		Anime:              anime.New(tx),
		Card:               card.New(tx),
		CardTag:            card_tag.New(tx),
		AccountCard:        account_card.New(tx),
		Pack:               pack.New(tx),
		AccountPackPity:    account_pack_pity.New(tx),
		AccountPullHistory: account_pull_history.New(tx),
		PackPool:           pack_pool.New(tx),
	}
}
