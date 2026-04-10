package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/card_tag"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type CardTagDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func CardTag(router *gin.RouterGroup, d CardTagDependency) {
	handler := card_tag.New(card_tag.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("card-tag")
	admin := r.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.GET("", handler.FindAll)
		admin.GET(":id", handler.FindOneById)
		admin.POST("", handler.CreateOne)
		admin.PATCH(":id", handler.UpdateOneById)
		admin.DELETE(":id", handler.DeleteOneById)
	}
}
