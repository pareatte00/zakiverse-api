package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type Dependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Router     *gin.Engine
	Middleware *middleware.Middleware
	Service    *service.Service
}

func Bind(d Dependency) {
	router := d.Router

	router.Use(d.Middleware.GetLocale)

	router.NoRoute(empty)
	router.GET("health", health)
	router.GET("info", info)

	v1 := router.Group("v1")
	v1.Use(d.Middleware.AuthSystemServiceKey)

	Account(v1, AccountDependency{
		Config:     d.Config,
		Credential: d.Credential,
		Middleware: d.Middleware,
		Service:    d.Service,
	})
}
