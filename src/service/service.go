package service

import (
	"database/sql"

	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/outbound"
	"github.com/zakiverse/zakiverse-api/src/repository"
)

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Database   *sql.DB
	Repository *repository.Repository
	Outbound   *outbound.Outbound
}

type Service struct {
	config     config.ConfigConstant
	credential config.ConfigCredential
	database   *sql.DB
	repository *repository.Repository
	outbound   *outbound.Outbound

	Account        *AccountService
	AccountBalance *AccountBalanceService
	Anime          *AnimeService
	Card           *CardService
	AccountCard    *AccountCardService
	CardTag        *CardTagService
	CheckIn        *CheckInService
	Pack           *PackService
	PackPool       *PackPoolService
	Profile        *ProfileService
}

func New(d Dependency) *Service {
	service := &Service{
		config:     d.Config,
		credential: d.Credential,
		database:   d.Database,
		repository: d.Repository,
		outbound:   d.Outbound,
	}

	service.Account = &AccountService{service: service}
	service.AccountBalance = &AccountBalanceService{service: service}
	service.Anime = &AnimeService{service: service}
	service.Card = &CardService{service: service}
	service.CardTag = &CardTagService{service: service}
	service.AccountCard = &AccountCardService{service: service}
	service.CheckIn = &CheckInService{service: service}
	service.Pack = &PackService{service: service}
	service.PackPool = &PackPoolService{service: service}
	service.Profile = &ProfileService{service: service}

	return service
}
