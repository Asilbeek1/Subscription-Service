package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Asilbeek1/Subscription-Service/internal/config"
	"github.com/Asilbeek1/Subscription-Service/internal/database"
	"github.com/Asilbeek1/Subscription-Service/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found. Copy .env.example to .env")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Init Config
	cfg := config.MustLoad()
	fmt.Println("Initialized Config")

	//Init Logger
	log := logger.SetUpLogger(cfg.Env)
	log.Info("Starting Subscription Service logger ")

	//Init db
	_, err := database.OpenDB(ctx, cfg.Postgres, log)
	if err != nil {
		log.Error("Database connection error")
		os.Exit(1)
	}
	fmt.Println("Initialized Database")

	//Init Server

}
