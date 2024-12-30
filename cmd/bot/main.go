package main

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/bot"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log"
)

func main() {
	app, err := config.LoadApp()
	if err != nil {
		log.Fatal(err)
	}

	sl.SetupLogger(app.Vars.Environment)

	if err = postgres.OpenDB(app.Vars); err != nil {
		log.Fatal(err)
	}

	err = translator.InitTranslator("./locales")
	if err != nil {
		log.Fatal(err)
	}

	if err = bot.Run(app); err != nil {
		log.Fatal(err)
	}
}
