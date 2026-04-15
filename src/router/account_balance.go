package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/account_balance"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type AccountBalanceDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func AccountBalance(router *gin.RouterGroup, d AccountBalanceDependency) {
	handler := account_balance.New(account_balance.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("account/balance")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("", handler.GetBalance)
	}
}
