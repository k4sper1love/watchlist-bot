package admin

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

// HandleFeedbackDetailCommand handles the command for viewing detailed information about a specific feedback.
// Retrieves the feedback details and sends a message with its information and an appropriate keyboard.
func HandleFeedbackDetailCommand(app models.App, session *models.Session) {
	if feedback, err := postgres.GetFeedbackByID(session.AdminState.FeedbackID); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallAdminFeedback))
	} else {
		app.SendMessage(messages.FeedbackDetail(session, feedback), keyboards.FeedbackDetail(session))
	}
}

// HandleFeedbackDetailButtons handles button interactions related to the feedback detail view.
// Supports actions like going back or deleting the feedback.
func HandleFeedbackDetailButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFeedbackDetailBack:
		general.RequireRole(app, session, HandleFeedbacksCommand, roles.Admin)

	case states.CallFeedbackDetailDelete:
		general.RequireRole(app, session, handleFeedbackDetailDelete, roles.Admin)
	}
}

// handleFeedbackDetailDelete processes the deletion of a specific feedback.
// Deletes the feedback from the database and navigates back to the feedback list.
func handleFeedbackDetailDelete(app models.App, session *models.Session) {
	if err := postgres.DeleteFeedbackByID(session.AdminState.FeedbackID); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallAdminFeedback))
		return
	}

	app.SendMessage(messages.FeedbackDeleteSuccess(session), nil)
	HandleFeedbacksCommand(app, session)
}
