package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"os"
	"strconv"
)

func LoadApp() (*models.App, error) {
	if err := godotenv.Load(); err != nil {
		sl.Log.Warn("not found .env file")
	}

	rootID, err := strconv.Atoi(os.Getenv("ROOT_TELEGRAM_ID"))
	if err != nil {
		sl.Log.Error("failed to parse root_telegram_id")
		return nil, err
	}

	vars := &models.Vars{
		Version:           os.Getenv("VERSION"),
		BotToken:          os.Getenv("BOT_TOKEN"),
		Environment:       os.Getenv("ENVIRONMENT"),
		DSN:               configureDSN(),
		Host:              os.Getenv("API_HOST"),
		Secret:            os.Getenv("API_SECRET"),
		RootID:            rootID,
		KinopoiskAPIToken: os.Getenv("KINOPOISK_API_TOKEN"),
		YoutubeAPIToken:   os.Getenv("YOUTUBE_API_TOKEN"),
		IMDBAPIToken:      os.Getenv("IMDB_API_TOKEN"),
	}
	sl.Log.Debug("vars successfully set")

	app := &models.App{
		Vars: vars,
	}
	sl.Log.Debug("app successfully created")

	return app, nil
}

func configureDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}
