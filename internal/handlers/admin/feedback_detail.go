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
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleFeedbackDetailCommand(app models.App, session *models.Session) {
	feedback, err := postgres.GetFeedbackByID(session.AdminState.FeedbackID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	msg := messages.BuildFeedbackDetailMessage(session, feedback)
	keyboard := keyboards.BuildAdminFeedbackDetailKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleFeedbackDetailButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackAdminFeedbackDetailBack:
		general.RequireRole(app, session, HandleFeedbacksCommand, roles.Admin)

	case callback == states.CallbackAdminFeedbackDetailDelete:
		general.RequireRole(app, session, handleFeedbackDetailDelete, roles.Admin)
	}
}

func handleFeedbackDetailDelete(app models.App, session *models.Session) {
	feedback, err := postgres.GetFeedbackByID(session.AdminState.FeedbackID)
	if err != nil {
		handleFeedbackDetailError(app, session)
		return
	}

	err = postgres.DeleteFeedbackByID(int(feedback.ID))
	if err != nil {
		handleFeedbackDetailError(app, session)
		return
	}

	msg := translator.Translate(session.Lang, "success", nil, nil)
	app.SendMessage(msg, nil)
	HandleFeedbacksCommand(app, session)
}

func handleFeedbackDetailError(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "someError", nil, nil)
	app.SendMessage(msg, nil)
	HandleFeedbackDetailCommand(app, session)
}
