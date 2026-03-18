package config

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/logger"
)

func ServerConfig(router *gin.Engine, port string) *http.Server {
	router.MaxMultipartMemory = 5 << 20

	s := &http.Server{
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if port == "" {
		log.Fatal("Failed to run server because port config is undefined")
	} else {
		s.Addr = ":" + port
		logger.Info("Server Running", logger.Field("port", port))
	}

	return s
}
