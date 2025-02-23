package validator

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func sendValidationMessage(app models.App, session *models.Session, messageID string, data map[string]interface{}) {
	msg := fmt.Sprintf("⚠️ %s", translator.Translate(session.Lang, messageID, data, nil))
	app.SendMessage(msg, nil)
}

func HandleInvalidInputLength(app models.App, session *models.Session, minLength, maxLength int) {
	sendValidationMessage(app, session, "invalidInputLength", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	})
}

func HandleInvalidInputRange[T int | float64](app models.App, session *models.Session, minValue T, maxValue T) {
	sendValidationMessage(app, session, "invalidInputRange", map[string]interface{}{
		"Min": minValue,
		"Max": maxValue,
	})
}

func HandleInvalidInputURL(app models.App, session *models.Session) {
	sendValidationMessage(app, session, "invalidInputURL", nil)
}
