package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	BotToken    string
	Environment string
	DSN         string
	BaseURL     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		Environment: os.Getenv("ENVIRONMENT"),
		DSN:         configureDSN(),
		BaseURL:     os.Getenv("BASE_URL"),
	}

	return cfg, nil
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
