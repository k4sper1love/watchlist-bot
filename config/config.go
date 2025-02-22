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

func InitAppConfig() (*models.App, error) {
	if err := godotenv.Load(); err != nil {
		sl.Log.Warn(".env file not found")
	}

	rootID, err := parseRootTelegramID()
	if err != nil {
		return nil, err
	}

	config := &models.Config{
		Version:         getEnvOrDefault("VERSION", "unknown"),
		BotToken:        getEnvOrDefault("BOT_TOKEN", ""),
		Environment:     getEnvOrDefault("ENVIRONMENT", "dev"),
		DatabaseURL:     buildDatabaseURL(),
		APIHost:         getEnvOrDefault("API_HOST", "localhost"),
		APISecret:       getEnvOrDefault("API_SECRET", ""),
		RootID:          rootID,
		YoutubeAPIToken: getEnvOrDefault("YOUTUBE_API_TOKEN", ""),
		IMDBAPIToken:    getEnvOrDefault("IMDB_API_TOKEN", ""),
	}

	return &models.App{Config: config}, nil
}

func buildDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnvOrDefault("POSTGRES_USER", "user"),
		getEnvOrDefault("POSTGRES_PASSWORD", "password"),
		getEnvOrDefault("POSTGRES_HOST", "localhost"),
		getEnvOrDefault("POSTGRES_PORT", "5432"),
		getEnvOrDefault("POSTGRES_DB", "database"),
	)
}

func parseRootTelegramID() (int, error) {
	idStr := os.Getenv("ROOT_TELEGRAM_ID")
	if idStr == "" {
		sl.Log.Error("ROOT_TELEGRAM_ID is not set")
		return 0, fmt.Errorf("ROOT_TELEGRAM_ID is not set")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sl.Log.Error("failed to parse ROOT_TELEGRAM_ID", slog.Any("error", err))
		return 0, fmt.Errorf("invalid ROOT_TELEGRAM_ID: %w", err)
	}
	return id, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		sl.Log.Warn("environment variable is not set, using default", slog.String("key", key), slog.String("default", defaultValue))
		return defaultValue
	}
	return value
}
