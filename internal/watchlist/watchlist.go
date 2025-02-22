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

	app, err := config.InitAppConfig()
	if err != nil {
		return err
	}
	sl.Log.Info("application config loaded successfully")

	sl.SetupLogger(app.Config.Environment)

	if err = postgres.ConnectDatabase(app.Config); err != nil {
		return err
	}
	sl.Log.Info("database connection established successfully")

	err = translator.InitTranslator("./locales")
	if err != nil {
		return err
	}
	sl.Log.Info("translator initialized successfully")

	return bot.Start(app)
}
