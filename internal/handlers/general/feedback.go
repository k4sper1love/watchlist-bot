package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"unicode/utf8"
)

func HandleFeedbackCommand(app models.App, session *models.Session) {
	msg := messages.BuildFeedbackMessage(session)

	keyboard := keyboards.BuildFeedbackKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleFeedbackButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackFeedbackCategorySuggestions:
		session.FeedbackState.Category = "offers"
	case states.CallbackFeedbackCategoryBugs:
		session.FeedbackState.Category = "mistakes"
	case states.CallbackFeedbackCategoryOther:
		session.FeedbackState.Category = "otherIssues"
	}

	handleFeedbackMessage(app, session)
}

func HandleFeedbackProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearState()
		HandleFeedbackCommand(app, session)
	}

	switch session.State {
	case states.ProcessFeedbackAwaitingMessage:
		parseFeedbackMessage(app, session)
	}
}

func handleFeedbackMessage(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "category", nil, nil)
	part2 := translator.Translate(session.Lang, session.FeedbackState.Category, nil, nil)
	part3 := translator.Translate(session.Lang, "feedbackTextRequest", nil, nil)

	msg := fmt.Sprintf("📄 <b>%s:</b> <code>%s</code>\n\n%s", part1, part2, part3)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessFeedbackAwaitingMessage)
}

func parseFeedbackMessage(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddBack("").Build(session.Lang)

	text := utils.ParseMessageString(app.Update)

	if utf8.RuneCountInString(text) > 3000 {
		part := translator.Translate(session.Lang, "maxLengthInSymbols", map[string]interface{}{
			"Length": 3000,
		}, nil)
		msg := fmt.Sprintf("⚠️ %s", part)
		app.SendMessage(msg, nil)

		handleFeedbackMessage(app, session)
		return
	}

	session.FeedbackState.Message = text

	err := postgres.SaveFeedback(session.TelegramID, session.TelegramUsername, session.FeedbackState.Category, session.FeedbackState.Message)
	if err != nil {
		part1 := translator.Translate(session.Lang, "feedbackFailure", nil, nil)
		part2 := translator.Translate(session.Lang, "tryLater", nil, nil)
		msg := fmt.Sprintf("🚨 %s\n%s", part1, part2)
		app.SendMessage(msg, keyboard)
		return
	}

	successMsg := translator.Translate(session.Lang, "feedbackSuccess", nil, nil)
	msg := fmt.Sprintf("✅ %s", successMsg)

	app.SendMessage(msg, keyboard)

	session.ClearAllStates()
}
