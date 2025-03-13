package main

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/bot"
	"os"
)

func main() {
	if err := bot.Run(); err != nil {
		sl.Log.Error("application terminated due to an error")
		os.Exit(1)
	}
}
