// Package main initializes and runs the Watchlist bot application.
//
// It serves as the entry point for the bot, handling initialization and error logging.
package main

import (
	"github.com/k4sper1love/watchlist-bot/internal/bot"
	"log/slog"
	"os"
)

// main is the entry point of the application.
// It initializes and runs the bot, logging any critical errors that cause termination.
func main() {
	if err := bot.Run(); err != nil {
		slog.Error("application terminated due to an error")
		os.Exit(1)
	}
}
