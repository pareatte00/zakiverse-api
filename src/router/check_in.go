package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/handler/check_in"
	"github.com/zakiverse/zakiverse-api/src/handler/check_in_plan"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/service"
)

type CheckInDependency struct {
	Config     config.ConfigConstant
	Credential config.ConfigCredential
	Middleware *middleware.Middleware
	Service    *service.Service
}

func CheckIn(router *gin.RouterGroup, d CheckInDependency) {
	playerHandler := check_in.New(check_in.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	adminHandler := check_in_plan.New(check_in_plan.Dependency{
		Config:     d.Config,
		Credential: d.Credential,
		Service:    d.Service,
	})

	// Player routes
	r := router.Group("check-in")
	auth := r.Use(d.Middleware.AuthJWT)
	{
		auth.GET("", playerHandler.GetPlans)
		auth.POST(":planId", playerHandler.Claim)
	}

	// Admin routes
	ar := router.Group("check-in-plan")
	admin := ar.Use(d.Middleware.AuthJWT, d.Middleware.AuthAdmin)
	{
		admin.GET("", adminHandler.FindAll)
		admin.GET(":id", adminHandler.FindOneById)
		admin.POST("", adminHandler.CreateOne)
		admin.PATCH(":id", adminHandler.UpdateOneById)
		admin.DELETE(":id", adminHandler.DeleteOneById)
	}
}
