package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"os"
)

func LoadApp() (*models.App, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	vars := &models.Vars{
		BotToken:    os.Getenv("BOT_TOKEN"),
		Environment: os.Getenv("ENVIRONMENT"),
		DSN:         configureDSN(),
		BaseURL:     os.Getenv("BASE_URL"),
		Secret:      os.Getenv("SECRET"),
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
