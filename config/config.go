package config

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"os"
)

type App struct {
	Vars *Vars
	Bot  *tgbotapi.BotAPI
	Upd  *tgbotapi.Update
}

type Vars struct {
	BotToken    string
	Environment string
	DSN         string
	BaseURL     string
	Secret      string
}

func LoadApp() (*App, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	vars := &Vars{
		BotToken:    os.Getenv("BOT_TOKEN"),
		Environment: os.Getenv("ENVIRONMENT"),
		DSN:         configureDSN(),
		BaseURL:     os.Getenv("BASE_URL"),
		Secret:      os.Getenv("SECRET"),
	}

	app := &App{
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
