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
	admin, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleUsersCommand(app, session)
		return
	}

	msg := messages.BuildAdminUserDetailMessage(session, admin)
	keyboard := keyboards.BuildAdminDetailKeyboard(session, admin)

	app.SendMessage(msg, keyboard)
}

func HandleAdminDetailButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

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
	user, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	msg := ""

	if session.Role == roles.Root {
		if user.Role == roles.SuperAdmin {
			msg = translator.Translate(session.Lang, "alreadyMaxRole", nil, nil)
		} else {
			success, err := postgres.SetUserRole(user.TelegramID, user.Role.NextRole())
			if err != nil || !success {
				msg = translator.Translate(session.Lang, "someError", nil, nil)
			} else {
				msg = translator.Translate(session.Lang, "success", nil, nil)
			}
		}
	} else if session.Role == roles.SuperAdmin {
		if user.Role == roles.Admin {
			msg = translator.Translate(session.Lang, "alreadyMaxRole", nil, nil)
		} else if user.Role.HasAccess(roles.SuperAdmin) {
			msg = translator.Translate(session.Lang, "noAccess", nil, nil)
		} else {
			success, err := postgres.SetUserRole(user.TelegramID, user.Role.NextRole())
			if err != nil || !success {
				msg = translator.Translate(session.Lang, "someError", nil, nil)
			} else {
				msg = translator.Translate(session.Lang, "success", nil, nil)
			}
		}
	} else {
		msg = translator.Translate(session.Lang, "noAccess", nil, nil)
	}

	app.SendMessage(msg, nil)

	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}

func handleLowerRole(app models.App, session *models.Session) {
	user, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	msg := ""

	if session.Role == roles.Root {
		if user.Role == roles.User {
			msg = translator.Translate(session.Lang, "alreadyMinRole", nil, nil)
		} else {
			success, err := postgres.SetUserRole(user.TelegramID, user.Role.PrevRole())
			if err != nil || !success {
				msg = translator.Translate(session.Lang, "someError", nil, nil)
			} else {
				msg = translator.Translate(session.Lang, "success", nil, nil)
			}
		}
	} else if session.Role == roles.SuperAdmin {
		if user.Role == roles.SuperAdmin || user.Role == roles.User {
			msg = translator.Translate(session.Lang, "noAccess", nil, nil)
		} else {
			success, err := postgres.SetUserRole(user.TelegramID, user.Role.PrevRole())
			if err != nil || !success {
				msg = translator.Translate(session.Lang, "someError", nil, nil)
			} else {
				msg = translator.Translate(session.Lang, "success", nil, nil)
			}
		}
	} else {
		msg = translator.Translate(session.Lang, "noAccess", nil, nil)
	}

	app.SendMessage(msg, nil)

	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}

func handleRemoveRole(app models.App, session *models.Session) {
	user, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
		return
	}

	msg := ""

	if !user.Role.HasAccess(session.Role) {
		_, err = postgres.SetUserRole(user.TelegramID, roles.User)
		if err != nil {
			msg = translator.Translate(session.Lang, "someError", nil, nil)
		} else {
			msg = translator.Translate(session.Lang, "success", nil, nil)
		}
	}

	app.SendMessage(msg, nil)

	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}
