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

func HandleAdminDetailCommand(app models.App, session *models.Session) {
	admin, err := postgres.GetUserByField(postgres.TelegramIDField, session.AdminState.UserID, true)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminSelectUsers).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		session.ClearAllStates()
		return
	}

	session.AdminState.UserLang = admin.Lang
	session.AdminState.UserRole = admin.Role

	msg := messages.BuildAdminUserDetailMessage(session, admin)
	keyboard := keyboards.BuildAdminDetailKeyboard(session, admin)

	app.SendMessage(msg, keyboard)
}

func HandleAdminDetailButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch {
	case callback == states.CallbackAdminDetailBack:
		general.RequireRole(app, session, HandleAdminsCommand, roles.Admin)

	case callback == states.CallbackAdminDetailRaiseRole:
		general.RequireRole(app, session, handleRaiseRole, roles.Admin)

	case callback == states.CallbackAdminDetailLowerRole:
		general.RequireRole(app, session, handleLowerRole, roles.Admin)

	case callback == states.CallbackAdminDetailRemoveRole:
		general.RequireRole(app, session, handleRemoveRole, roles.SuperAdmin)

	}
}

func handleRaiseRole(app models.App, session *models.Session) {
	var msg string

	if !session.Role.HasAccess(roles.SuperAdmin) || (session.AdminState.UserRole == roles.Admin && session.Role != roles.Root) {
		msg = "❗️" + translator.Translate(session.Lang, "noAccess", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	if session.AdminState.UserRole == roles.SuperAdmin {
		msg = "❗️" + translator.Translate(session.Lang, "alreadyMaxRole", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	err := postgres.SetUserRole(session.AdminState.UserID, session.AdminState.UserRole.NextRole())
	if err != nil {
		msg = "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
	} else {
		msg = messages.BuildRaiseRoleNotificationMessage(session)
		app.SendMessageByID(session.AdminState.UserID, msg, nil)

		msg = messages.BuildRaiseRoleMessage(session)
	}

	app.SendMessage(msg, nil)

	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}

func handleLowerRole(app models.App, session *models.Session) {
	msg := ""

	if session.AdminState.UserID == session.TelegramID {
		msg = "❗️" + translator.Translate(session.Lang, "cannotLowerSelf", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	if !session.Role.HasAccess(roles.SuperAdmin) || session.AdminState.UserRole.HasAccess(roles.SuperAdmin) && !session.Role.HasAccess(roles.Root) {
		msg = "❗️" + translator.Translate(session.Lang, "noAccess", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	if session.AdminState.UserRole == roles.User {
		msg = "❗️" + translator.Translate(session.Lang, "alreadyMinRole", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	err := postgres.SetUserRole(session.AdminState.UserID, session.AdminState.UserRole.PrevRole())
	if err != nil {
		msg = "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
	} else {
		msg = messages.BuildLowerRoleNotificationMessage(session)
		app.SendMessageByID(session.AdminState.UserID, msg, nil)

		msg = messages.BuildLowerRoleMessage(session)
	}

	app.SendMessage(msg, nil)
	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}

func handleRemoveRole(app models.App, session *models.Session) {
	msg := ""

	if !session.AdminState.UserRole.HasAccess(session.Role) {
		err := postgres.SetUserRole(session.AdminState.UserID, roles.User)
		if err != nil {
			msg = "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		} else {
			msg = messages.BuildRemoveRoleNotificationMessage(session)
			app.SendMessageByID(session.AdminState.UserID, msg, nil)

			msg = messages.BuildRemoveRoleMessage(session)
		}
	}

	app.SendMessage(msg, nil)

	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}
