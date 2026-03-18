package config

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func SetTimeout(router *gin.Engine, max int) {
	router.Use(timeout.New(
		timeout.WithTimeout(time.Duration(max)*time.Second),
		timeout.WithResponse(func(c *gin.Context) {
			c.AbortWithStatus(http.StatusRequestTimeout)
		}),
	))
}
