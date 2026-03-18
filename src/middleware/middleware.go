package middleware

import (
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type Middleware struct {
	config     config.ConfigConstant
	credential config.ConfigCredential
	service    *service.Service
}

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Service    *service.Service
}

func New(d Dependency) *Middleware {
	return &Middleware{
		config:     d.Config,
		credential: d.Credential,
		service:    d.Service,
	}
}
