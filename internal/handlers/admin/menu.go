package admin

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

func HandleMenuCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildAdminMenuMessage(session), keyboards.BuildAdminMenuKeyboard(session))
}

func HandleMenuButton(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackAdminSelectAdmins:
		session.AdminState.IsAdmin = true
		resetAdminPageAndHandle(app, session, HandleEntitiesCommand, roles.Admin)

	case states.CallbackAdminSelectUsers:
		session.AdminState.IsAdmin = false
		resetAdminPageAndHandle(app, session, HandleEntitiesCommand, roles.Admin)

	case states.CallbackAdminSelectBroadcast:
		resetAdminPageAndHandle(app, session, HandleBroadcastCommand, roles.Admin)

	case states.CallbackAdminSelectFeedback:
		resetAdminPageAndHandle(app, session, HandleFeedbacksCommand, roles.Helper)
	}
}

func resetAdminPageAndHandle(app models.App, session *models.Session, next func(models.App, *models.Session), role roles.Role) {
	session.AdminState.CurrentPage = 1
	general.RequireRole(app, session, next, role)
}
