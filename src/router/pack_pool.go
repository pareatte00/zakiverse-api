package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/pack_pool"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type PackPoolDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func PackPool(router *gin.RouterGroup, d PackPoolDependency) {
	handler := pack_pool.New(pack_pool.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("pack-pool")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("active", handler.FindActiveBanners)
		auth.GET(":id", handler.FindOneWithPacks)
	}
	admin := r.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.GET("", handler.FindAll)
		admin.POST("", handler.CreateOne)
		admin.POST("reorder", handler.Reorder)
		admin.GET(":id/detail", handler.FindOneById)
		admin.PATCH(":id", handler.UpdateOneById)
		admin.DELETE(":id", handler.DeleteOneById)
		admin.POST(":id/reorder-packs", handler.ReorderPacks)
	}
}
