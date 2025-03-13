package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

func HandleFeedbackCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Feedback(session), keyboards.Feedback(session))
}

func HandleFeedbackButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	default:
		if strings.HasPrefix(callback, states.FeedbackCategory) {
			session.FeedbackState.Category = strings.TrimPrefix(callback, states.FeedbackCategory)
			requestFeedbackMessage(app, session)
		}
	}
}

func HandleFeedbackProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearState()
		HandleFeedbackCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitFeedbackMessage:
		parser.ParseFeedbackMessage(app, session, requestFeedbackMessage, saveFeedback)
	}
}

func requestFeedbackMessage(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFeedbackMessage(session), keyboards.Cancel(session))
	session.SetState(states.AwaitFeedbackMessage)
}

func saveFeedback(app models.App, session *models.Session) {
	if err := postgres.SaveFeedback(session.TelegramID, session.TelegramUsername, session.FeedbackState.Category, session.FeedbackState.Message); err != nil {
		app.SendMessage(messages.FeedbackFailure(session), keyboards.Back(session, ""))
	} else {
		app.SendMessage(messages.FeedbackSuccess(session), keyboards.Back(session, ""))
	}

	session.ClearAllStates()
}
