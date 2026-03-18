package config

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
)

func SetRateLimit(router *gin.Engine, limit int) {
	limiter := tollbooth.NewLimiter(float64(limit), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	router.Use(rateLimitHandler(limiter))
}

func rateLimitHandler(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		} else {
			c.Next()
		}
	}
}
