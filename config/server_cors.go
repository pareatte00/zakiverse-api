package config

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetCors(router *gin.Engine, allowDomain []string) {
	config := cors.DefaultConfig()
	config.AllowOrigins = allowDomain
	config.AllowMethods = []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
	}
	config.AllowHeaders = []string{
		"Origin", "Content-Length", "Content-Type", "Accept-Language", "Authorization",
		"X-Locale", "X-System-Key",
		"X-Actor-Id", "X-Actor-Type",
	}
	config.AllowCredentials = true
	config.MaxAge = 10 * time.Minute

	router.Use(cors.New(config))
}
