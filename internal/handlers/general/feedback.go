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
)

func HandleFeedbackCommand(app models.App, session *models.Session) {
	msg := messages.BuildFeedbackMessage(session)

	keyboard := keyboards.BuildFeedbackKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleFeedbackButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	var category string
	switch callback {
	case states.CallbackFeedbackCategorySuggestions:
		category = "Offers"
	case states.CallbackFeedbackCategoryBugs:
		category = "Mistakes"
	case states.CallbackFeedbackCategoryOther:
		category = "Other issues"
	}

	session.FeedbackState.Category = category

	part1 := translator.Translate(session.Lang, "feedbackCurrentCategory", nil, nil)
	part2 := translator.Translate(session.Lang, "feedbackTextRequest", nil, nil)

	msg := fmt.Sprintf("📄 <b>%s:</b> %s\n\n%s", part1, category, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessFeedbackAwaitingMessage)
}

func HandleFeedbackProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		HandleFeedbackCommand(app, session)
	}

	switch session.State {
	case states.ProcessFeedbackAwaitingMessage:
		parseFeedbackMessage(app, session)
	}
}

func parseFeedbackMessage(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddBack("").Build(session.Lang)

	session.FeedbackState.Message = utils.ParseMessageString(app.Upd)

	err := postgres.SaveFeedbackToDatabase(session.TelegramID, session.FeedbackState.Category, session.FeedbackState.Message)
	if err != nil {
		part1 := translator.Translate(session.Lang, "feedbackFailure", nil, nil)
		part2 := translator.Translate(session.Lang, "tryLater", nil, nil)
		msg := fmt.Sprintf("❌%s\n%s", part1, part2)
		app.SendMessage(msg, keyboard)
		return
	}

	successMsg := translator.Translate(session.Lang, "feedbackSuccess", nil, nil)
	msg := fmt.Sprintf("✅ %s", successMsg)

	app.SendMessage(msg, keyboard)

	session.ClearAllStates()
}
