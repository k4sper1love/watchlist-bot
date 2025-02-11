package watchlist

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/bot"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func Run() error {
	sl.SetupLogger("dev")
	sl.Log.Info("starting application")

	sl.Log.Info("loading application config...")
	app, err := config.LoadApp()
	if err != nil {
		return err
	}

	sl.SetupLogger(app.Vars.Environment)

	sl.Log.Info("opening database connection...")
	if err = postgres.OpenDB(app.Vars); err != nil {
		return err
	}

	sl.Log.Info("initializing translator...")
	err = translator.InitTranslator("./locales")
	if err != nil {
		return err
	}

	sl.Log.Info("starting bot...")
	return bot.Run(app)
}
