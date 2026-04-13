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
		auth.GET(":id/next-packs", handler.FindNextPacks)
	}
	admin := r.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.GET("", handler.FindAll)
		admin.POST("", handler.CreateOne)
		admin.POST("sort", handler.Sort)
		admin.PATCH(":id", handler.UpdateOneById)
		admin.DELETE(":id", handler.DeleteOneById)
		admin.POST(":id/assign-packs", handler.AssignPacks)
		admin.POST(":id/sort-packs", handler.SortPacks)
		admin.POST(":id/sort-rotation", handler.SortRotation)
	}
}
