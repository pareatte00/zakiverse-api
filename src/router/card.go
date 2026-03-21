package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/card"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type CardDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func Card(router *gin.RouterGroup, d CardDependency) {
	handler := card.New(card.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("card")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("anime/:animeId", handler.FindAllByAnimeId)
		auth.GET(":id", handler.FindOneById)
	}
	admin := r.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.POST("", handler.CreateOne)
		admin.PUT(":id", handler.UpdateOneById)
		admin.DELETE(":id", handler.DeleteOneById)
	}
}
