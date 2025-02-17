package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
	"os"
	"strconv"
)

func LoadApp() (*models.App, error) {
	if err := godotenv.Load(); err != nil {
		sl.Log.Warn("not found .env file")
	}

	rootID, err := strconv.Atoi(os.Getenv("ROOT_TELEGRAM_ID"))
	if err != nil {
		sl.Log.Error("failed to parse root_telegram_id", slog.Any("error", err))
		return nil, err
	}

	vars := &models.Vars{
		Version:         os.Getenv("VERSION"),
		BotToken:        os.Getenv("BOT_TOKEN"),
		Environment:     os.Getenv("ENVIRONMENT"),
		DSN:             configureDSN(),
		Host:            os.Getenv("API_HOST"),
		Secret:          os.Getenv("API_SECRET"),
		RootID:          rootID,
		YoutubeAPIToken: os.Getenv("YOUTUBE_API_TOKEN"),
		IMDBAPIToken:    os.Getenv("IMDB_API_TOKEN"),
	}

	app := &models.App{
		Vars: vars,
	}

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
