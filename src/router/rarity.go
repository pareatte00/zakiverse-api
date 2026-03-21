package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/rarity"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type RarityDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func Rarity(router *gin.RouterGroup, d RarityDependency) {
	handler := rarity.New(rarity.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("rarity")
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
