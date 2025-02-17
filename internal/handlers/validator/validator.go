package validator

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleInvalidInputLength(app models.App, session *models.Session, minLength, maxLength int) {
	part := translator.Translate(session.Lang, "invalidInputLength", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	}, nil)
	msg := fmt.Sprintf("⚠️ %s", part)
	app.SendMessage(msg, nil)
}

func HandleInvalidInputRange[T int | float64](app models.App, session *models.Session, minValue T, maxValue T) {
	part := translator.Translate(session.Lang, "invalidInputRange", map[string]interface{}{
		"Min": minValue,
		"Max": maxValue,
	}, nil)
	msg := fmt.Sprintf("⚠️ %s", part)
	app.SendMessage(msg, nil)
}

func HandleInvalidInputURL(app models.App, session *models.Session) {
	part := translator.Translate(session.Lang, "invalidInputURL", nil, nil)
	msg := fmt.Sprintf("⚠️ %s", part)
	app.SendMessage(msg, nil)
}
