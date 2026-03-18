package repository

import (
	"database/sql"

	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/repository/account"
)

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Database   *sql.DB
}

type Repository struct {
	config     config.ConfigConstant
	credential config.ConfigCredential

	Account *account.Repository
}

func New(d Dependency) *Repository {
	return &Repository{
		config:     d.Config,
		credential: d.Credential,

		Account: account.New(d.Database),
	}
}

func (r *Repository) Tx(tx *sql.Tx) *Repository {
	return &Repository{
		config:     r.config,
		credential: r.credential,

		Account: account.New(tx),
	}
}
