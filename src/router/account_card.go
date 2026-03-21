package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	accountCardHandler "github.com/zakiverse/zakiverse-api/src/handler/account_card"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type AccountCardDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func AccountCard(router *gin.RouterGroup, d AccountCardDependency) {
	handler := accountCardHandler.New(accountCardHandler.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("account-card")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("me", handler.FindMyCards)
	}
	admin := r.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.POST("", handler.AddCard)
		admin.DELETE("", handler.RemoveCard)
	}
}
