package main

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/bot"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	sl.SetupLogger(cfg.Environment)

	if err = postgres.OpenDB(cfg); err != nil {
		log.Fatal(err)
	}

	if err = bot.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
