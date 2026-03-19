package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/config"
	cst "github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/database"
	"github.com/zakiverse/zakiverse-api/logger"
	"github.com/zakiverse/zakiverse-api/src/middleware"
	"github.com/zakiverse/zakiverse-api/src/outbound"
	"github.com/zakiverse/zakiverse-api/src/repository"
	"github.com/zakiverse/zakiverse-api/src/router"
	"github.com/zakiverse/zakiverse-api/src/service"
)

func init() {
	log.Printf("Initializing services...")
}

func main() {
	// Config
	conf, credential := config.InitConfig()

	// Logger
	logger.InitLogger(conf.Application.DeployMode)

	// Database
	dbConn := database.InitDB(conf, credential)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ginRouter := gin.Default()
	config.SetCors(ginRouter, conf.Application.CorsAllowOrigins)
	config.SetRateLimit(ginRouter, conf.Application.RateLimitPerSecond)
	config.SetTimeout(ginRouter, conf.Application.Timeout)
	server := config.ServerConfig(ginRouter, strconv.Itoa(conf.Application.DeployPort))

	repo := repository.New(repository.Dependency{
		Config:     conf,
		Credential: credential,
		Database:   dbConn,
	})
	ob := outbound.New(conf)
	serv := service.New(service.Dependency{
		Config:     conf,
		Credential: credential,
		Database:   dbConn,
		Repository: repo,
		Outbound:   ob,
	})
	mw := middleware.New(middleware.Dependency{
		Config:     conf,
		Credential: credential,
		Service:    serv,
	})
	router.Bind(router.Dependency{
		Router:     ginRouter,
		Config:     conf,
		Credential: credential,
		Middleware: mw,
		Service:    serv,
	})

	logger.Info("Services are now running and on standby")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Service ListenAndServe fail", logger.Field(cst.KeyError, err))
		}
	}()

	<-ctx.Done()

	logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	defer func() {
		if err := dbConn.Close(); err != nil {
			logger.Error("Failed to close database connection", logger.Field(cst.KeyError, err))
		} else {
			logger.Info("Database connection closed successfully")
		}
	}()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Forced to shutdown fail", logger.Field(cst.KeyError, err))
	}

	logger.Info("Service shutting down completely")
}
