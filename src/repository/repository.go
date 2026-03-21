package repository

import (
	"database/sql"

	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/repository/account"
	"github.com/zakiverse/zakiverse-api/src/repository/account_card"
	"github.com/zakiverse/zakiverse-api/src/repository/anime"
	"github.com/zakiverse/zakiverse-api/src/repository/card"
	"github.com/zakiverse/zakiverse-api/src/repository/rarity"
)

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Database   *sql.DB
}

type Repository struct {
	config     config.ConfigConstant
	credential config.ConfigCredential

	Account     *account.Repository
	Rarity      *rarity.Repository
	Anime       *anime.Repository
	Card        *card.Repository
	AccountCard *account_card.Repository
}

func New(d Dependency) *Repository {
	return &Repository{
		config:     d.Config,
		credential: d.Credential,

		Account:     account.New(d.Database),
		Rarity:      rarity.New(d.Database),
		Anime:       anime.New(d.Database),
		Card:        card.New(d.Database),
		AccountCard: account_card.New(d.Database),
	}
}

func (r *Repository) Tx(tx *sql.Tx) *Repository {
	return &Repository{
		config:     r.config,
		credential: r.credential,

		Account:     account.New(tx),
		Rarity:      rarity.New(tx),
		Anime:       anime.New(tx),
		Card:        card.New(tx),
		AccountCard: account_card.New(tx),
	}
}
