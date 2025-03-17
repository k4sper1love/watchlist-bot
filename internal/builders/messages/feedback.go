package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

// Feedback generates a message prompting the user to provide feedback.
func Feedback(session *models.Session) string {
	return fmt.Sprintf("💬 %s\n\n%s 😊\n\n%s",
		toBold(translator.Translate(session.Lang, "feedbackMessageHeader", nil, nil)),
		toItalic(translator.Translate(session.Lang, "feedbackMessageBody", nil, nil)),
		translator.Translate(session.Lang, "feedbackCategoryChoice", nil, nil))
}

// RequestFeedbackMessage generates a message prompting the user to enter feedback text.
// Includes the selected feedback category.
func RequestFeedbackMessage(session *models.Session) string {
	return fmt.Sprintf("📄 %s: %s\n\n%s",
		toBold(translator.Translate(session.Lang, "category", nil, nil)),
		toCode(translator.Translate(session.Lang, session.FeedbackState.Category, nil, nil)),
		translator.Translate(session.Lang, "feedbackTextRequest", nil, nil))
}

// FeedbackFailure generates an error message when submitting feedback fails.
func FeedbackFailure(session *models.Session) string {
	return fmt.Sprintf("🚨 %s\n%s",
		translator.Translate(session.Lang, "feedbackFailure", nil, nil),
		translator.Translate(session.Lang, "tryLater", nil, nil))
}

// FeedbackSuccess generates a success message after feedback is successfully submitted.
func FeedbackSuccess(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "feedbackSuccess", nil, nil)
}
