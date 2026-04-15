package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	profileHandler "github.com/zakiverse/zakiverse-api/src/handler/profile"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type ProfileDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func Profile(router *gin.RouterGroup, d ProfileDependency) {
	handler := profileHandler.New(profileHandler.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	r := router.Group("profile")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("me", handler.GetMe)
		auth.PATCH("me", handler.UpdateMe)
		auth.GET("search", handler.Search)
		auth.GET(":identifier", handler.GetByIdOrName)
	}
}
