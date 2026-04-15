package account_balance

import (
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type Handler struct {
	config     config.ConfigConstant
	credential config.ConfigCredential
	service    *service.Service
}

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Service    *service.Service
}

func New(d Dependency) Handler {
	return Handler{
		config:     d.Config,
		credential: d.Credential,
		service:    d.Service,
	}
}
