package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildFeedbackMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "feedbackMessageHeader", nil, nil)
	part2 := translator.Translate(session.Lang, "feedbackMessageBody", nil, nil)
	part3 := translator.Translate(session.Lang, "feedbackCategoryChoice", nil, nil)

	return fmt.Sprintf("💬 <b>%s</b>\n\n<i>%s</i> 😊\n\n%s", part1, part2, part3)
}

func BuildFeedbackRequestMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "category", nil, nil)
	part2 := translator.Translate(session.Lang, session.FeedbackState.Category, nil, nil)
	part3 := translator.Translate(session.Lang, "feedbackTextRequest", nil, nil)

	return fmt.Sprintf("📄 <b>%s:</b> <code>%s</code>\n\n%s", part1, part2, part3)
}

func BuildFeedbackMaxLengthMessage(session *models.Session, maxLength int) string {
	return "⚠️ " + translator.Translate(session.Lang, "maxLengthInSymbols", map[string]interface{}{
		"Length": maxLength,
	}, nil)
}

func BuildFeedbackFailureMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "feedbackFailure", nil, nil)
	part2 := translator.Translate(session.Lang, "tryLater", nil, nil)
	return fmt.Sprintf("🚨 %s\n%s", part1, part2)
}

func BuildFeedbackSuccessMessage(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "feedbackSuccess", nil, nil)
}
