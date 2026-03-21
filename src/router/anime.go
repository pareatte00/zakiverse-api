package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/anime"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type AnimeDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func Anime(router *gin.RouterGroup, d AnimeDependency) {
	handler := anime.New(anime.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("anime")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("", handler.FindAll)
		auth.GET(":id", handler.FindOneById)
	}
	admin := r.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.POST("", handler.CreateOne)
		admin.PUT(":id", handler.UpdateOneById)
		admin.DELETE(":id", handler.DeleteOneById)
	}
}
