package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
	"unicode/utf8"
)

const maxFeedbackLength = 3000

func HandleFeedbackCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Feedback(session), keyboards.Feedback(session))
}

func HandleFeedbackButtons(app models.App, session *models.Session) {
	switch {
	case strings.HasPrefix(utils.ParseCallback(app.Update), states.PrefixFeedbackCategory):
		session.FeedbackState.Category = strings.TrimPrefix(utils.ParseCallback(app.Update), states.PrefixFeedbackCategory)
		handleFeedbackMessage(app, session)
	}
}

func HandleFeedbackProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearState()
		HandleFeedbackCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFeedbackAwaitingMessage:
		parseFeedbackMessage(app, session)
	}
}

func handleFeedbackMessage(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFeedbackMessage(session), keyboards.Cancel(session))
	session.SetState(states.ProcessFeedbackAwaitingMessage)
}

func parseFeedbackMessage(app models.App, session *models.Session) {
	text := utils.ParseMessageString(app.Update)
	if utf8.RuneCountInString(text) > maxFeedbackLength {
		app.SendMessage(messages.WarningMaxLength(session, maxFeedbackLength), nil)
		handleFeedbackMessage(app, session)
		return
	}

	session.FeedbackState.Message = text
	saveFeedback(app, session)
	session.ClearAllStates()
}

func saveFeedback(app models.App, session *models.Session) {
	if err := postgres.SaveFeedback(session.TelegramID, session.TelegramUsername, session.FeedbackState.Category, session.FeedbackState.Message); err != nil {
		app.SendMessage(messages.FeedbackFailure(session), keyboards.Back(session, ""))
		return
	}

	app.SendMessage(messages.FeedbackSuccess(session), keyboards.Back(session, ""))
}
