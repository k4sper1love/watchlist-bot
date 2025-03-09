package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func Feedback(session *models.Session) string {
	return fmt.Sprintf("💬 %s\n\n%s 😊\n\n%s",
		toBold(translator.Translate(session.Lang, "feedbackMessageHeader", nil, nil)),
		toItalic(translator.Translate(session.Lang, "feedbackMessageBody", nil, nil)),
		translator.Translate(session.Lang, "feedbackCategoryChoice", nil, nil))
}

func RequestFeedbackMessage(session *models.Session) string {
	return fmt.Sprintf("📄 %s: %s\n\n%s",
		toBold(translator.Translate(session.Lang, "category", nil, nil)),
		toCode(translator.Translate(session.Lang, session.FeedbackState.Category, nil, nil)),
		translator.Translate(session.Lang, "feedbackTextRequest", nil, nil))
}

func WarningMaxLength(session *models.Session, maxLength int) string {
	return "⚠️ " + translator.Translate(session.Lang, "maxLengthInSymbols", map[string]interface{}{
		"Length": maxLength,
	}, nil)
}

func FeedbackFailure(session *models.Session) string {
	return fmt.Sprintf("🚨 %s\n%s",
		translator.Translate(session.Lang, "feedbackFailure", nil, nil),
		translator.Translate(session.Lang, "tryLater", nil, nil))
}

func FeedbackSuccess(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "feedbackSuccess", nil, nil)
}
