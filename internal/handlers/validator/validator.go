package validator

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func HandleInvalidInputLength(app models.App, session *models.Session, minLength, maxLength int) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputLength", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	}), nil)
}

func HandleInvalidInputRange[T int | float64](app models.App, session *models.Session, minValue T, maxValue T) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputRange", map[string]interface{}{
		"Min": minValue,
		"Max": maxValue,
	}), nil)
}

func HandleInvalidInputURL(app models.App, session *models.Session, minLength, maxLength int) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputURL", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	}), nil)
}

func HandleInvalidInputEmail(app models.App, session *models.Session, minLength, maxLength int) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputEmail", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	}), nil)
}
