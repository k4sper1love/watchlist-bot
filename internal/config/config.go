// Package config provides functionality for loading and managing application configuration.
//
// It reads environment variables, parses required values, and initializes the application's configuration.
package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
	"os"
	"strconv"
)

// Init initializes the application configuration by reading environment variables.
// It loads the .env file, parses required values (e.g., ROOT_TELEGRAM_ID), and sets defaults for missing values.
func Init() (*models.App, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found")
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
		LocalesDir:      getEnvOrDefault("LOCALES_DIR", "./locales"),
		APIHost:         getEnvOrDefault("API_HOST", "localhost"),
		APISecret:       getEnvOrDefault("API_SECRET", ""),
		RootID:          rootID,
		YoutubeAPIToken: getEnvOrDefault("YOUTUBE_API_TOKEN", ""),
		IMDBAPIToken:    getEnvOrDefault("IMDB_API_TOKEN", ""),
	}

	return &models.App{Config: config}, nil
}

// buildDatabaseURL constructs the PostgreSQL database connection URL from environment variables.
// It uses default values if any required environment variable is missing.
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

// parseRootTelegramID parses the ROOT_TELEGRAM_ID environment variable as an integer.
func parseRootTelegramID() (int, error) {
	idStr := os.Getenv("ROOT_TELEGRAM_ID")
	if idStr == "" {
		slog.Error("ROOT_TELEGRAM_ID is not set")
		return 0, fmt.Errorf("ROOT_TELEGRAM_ID is not set")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("failed to parse ROOT_TELEGRAM_ID", slog.Any("error", err))
		return 0, fmt.Errorf("invalid ROOT_TELEGRAM_ID: %w", err)
	}
	return id, nil
}

// getEnvOrDefault retrieves the value of an environment variable or returns a default value if it is not set.
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Warn("environment variable is not set, using default", slog.String("key", key), slog.String("default", defaultValue))
		return defaultValue
	}
	return value
}
