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

func HandleFeedbackDetailCommand(app models.App, session *models.Session) {
	if feedback, err := postgres.GetFeedbackByID(session.AdminState.FeedbackID); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallbackAdminSelectFeedback))
	} else {
		app.SendMessage(messages.FeedbackDetail(session, feedback), keyboards.FeedbackDetail(session))
	}
}

func HandleFeedbackDetailButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackAdminFeedbackDetailBack:
		general.RequireRole(app, session, HandleFeedbacksCommand, roles.Admin)

	case states.CallbackAdminFeedbackDetailDelete:
		general.RequireRole(app, session, handleFeedbackDetailDelete, roles.Admin)
	}
}

func handleFeedbackDetailDelete(app models.App, session *models.Session) {
	if err := postgres.DeleteFeedbackByID(session.AdminState.FeedbackID); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallbackAdminSelectFeedback))
		return
	}

	app.SendMessage(messages.FeedbackDeleteSuccess(session), nil)
	HandleFeedbacksCommand(app, session)
}
