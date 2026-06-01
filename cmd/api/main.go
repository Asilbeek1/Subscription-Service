package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Asilbeek1/Subscription-Service/internal/config"
	"github.com/Asilbeek1/Subscription-Service/internal/database"
	"github.com/Asilbeek1/Subscription-Service/internal/logger"
	"github.com/Asilbeek1/Subscription-Service/internal/service"
	"github.com/Asilbeek1/Subscription-Service/internal/transport/http/server"
)

// @title           Subscription Service API
// @version         1.0
// @description     Manages user subscriptions
// @host            localhost:8080
// @BasePath        /
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Init Config
	cfg := config.MustLoad()

	//Init Logger
	log := logger.SetUpLogger(cfg.Env)
	log.Info("Starting Subscription Service logger ")

	//Init database layer
	pool, err := database.OpenDB(ctx, cfg.Postgres, log)
	if err != nil {
		log.Error("Database connection error", "error", err)
		os.Exit(1)
	}
	log.Info("Database connection established", "port", cfg.Postgres.Port)

	//INIT Service
	service := service.NewSubscriptionService(pool, log)

	//Init Server
	router := server.New(service)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	log.Info("starting server", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("error running server", "error", err)
		os.Exit(1)
	}
}
