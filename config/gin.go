package config

import "github.com/gin-gonic/gin"

func setGinReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}