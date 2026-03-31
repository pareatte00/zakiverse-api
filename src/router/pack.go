package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/pack"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type PackDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func Pack(router *gin.RouterGroup, d PackDependency) {
	handler := pack.New(pack.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("pack")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("", handler.FindAll)
		auth.GET(":id", handler.FindOneById)
		auth.POST(":id/pull", handler.Pull)
	}
	admin := r.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.POST("", handler.CreateOne)
		admin.PATCH(":id", handler.UpdateOneById)
		admin.DELETE(":id", handler.DeleteOneById)
		admin.POST(":id/cards", handler.AddCards)
		admin.DELETE(":id/cards", handler.RemoveCards)
	}
}
