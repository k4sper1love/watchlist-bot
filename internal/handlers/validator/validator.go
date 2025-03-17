// Package validator handles input validation in the Watchlist application.
//
// It provides functions for checking constraints like length, range, URL, and email format,
// ensuring users receive clear feedback and guidance on correct input.
package validator

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

// HandleInvalidInputLength sends a validation warning to the user when the input length is outside the allowed range.
func HandleInvalidInputLength(app models.App, session *models.Session, minLength, maxLength int) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputLength", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	}), nil)
}

// HandleInvalidInputRange sends a validation warning to the user when the input value is outside the allowed range.
// This function supports both integer and float64 types for the range values.
func HandleInvalidInputRange[T int | float64](app models.App, session *models.Session, minValue T, maxValue T) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputRange", map[string]interface{}{
		"Min": minValue,
		"Max": maxValue,
	}), nil)
}

// HandleInvalidInputURL sends a validation warning to the user when the input URL is invalid or its length is outside the allowed range.
func HandleInvalidInputURL(app models.App, session *models.Session, minLength, maxLength int) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputURL", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	}), nil)
}

// HandleInvalidInputEmail sends a validation warning to the user when the input email is invalid or its length is outside the allowed range.
func HandleInvalidInputEmail(app models.App, session *models.Session, minLength, maxLength int) {
	app.SendMessage(messages.ValidationWarning(session, "invalidInputEmail", map[string]interface{}{
		"Min": minLength,
		"Max": maxLength,
	}), nil)
}
