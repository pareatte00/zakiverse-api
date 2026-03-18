package service

import (
	"database/sql"

	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/repository"
)

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Database   *sql.DB
	Repository *repository.Repository
}

type Service struct {
	config     config.ConfigConstant
	credential config.ConfigCredential
	database   *sql.DB
	repository *repository.Repository

	Account *AccountService
}

func New(d Dependency) *Service {
	service := &Service{
		config:     d.Config,
		credential: d.Credential,
		database:   d.Database,
		repository: d.Repository,
	}

	service.Account = &AccountService{service: service}

	return service
}
