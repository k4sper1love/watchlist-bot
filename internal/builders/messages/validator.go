package messages

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildValidationMessage(session *models.Session, messageID string, data map[string]interface{}) string {
	return "⚠️ " + translator.Translate(session.Lang, messageID, data, nil)
}
