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

// HandleFeedbackCommand handles the command for submitting feedback.
// Sends a message with options to select a feedback category and prompts the user to provide details.
func HandleFeedbackCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Feedback(session), keyboards.Feedback(session))
}

// HandleFeedbackButtons handles button interactions related to feedback submission.
// Supports actions like selecting a feedback category.
func HandleFeedbackButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	if strings.HasPrefix(callback, states.FeedbackCategory) {
		session.FeedbackState.Category = strings.TrimPrefix(callback, states.FeedbackCategory)
		requestFeedbackMessage(app, session)
	}
}

// HandleFeedbackProcess processes the workflow for submitting feedback.
// Handles states like awaiting input for the feedback message.
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

// requestFeedbackMessage prompts the user to enter a feedback message.
func requestFeedbackMessage(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFeedbackMessage(session), keyboards.Cancel(session))
	session.SetState(states.AwaitFeedbackMessage)
}

// saveFeedback saves the feedback to the database using the Postgres service.
func saveFeedback(app models.App, session *models.Session) {
	if err := postgres.SaveFeedback(session.TelegramID, session.TelegramUsername, session.FeedbackState.Category, session.FeedbackState.Message); err != nil {
		app.SendMessage(messages.FeedbackFailure(session), keyboards.Back(session, ""))
	} else {
		app.SendMessage(messages.FeedbackSuccess(session), keyboards.Back(session, ""))
	}

	session.ClearAllStates()
}
