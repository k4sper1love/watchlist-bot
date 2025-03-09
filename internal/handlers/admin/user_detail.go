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
	"strings"
)

func HandleUserDetailCommand(app models.App, session *models.Session) {
	if user, err := getEntity(session); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminSelectUsers))
	} else {
		app.SendMessage(messages.UserDetail(session, user), keyboards.BuildAdminUserDetailKeyboard(session, user))
	}
}

func HandleUserDetailButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackAdminUserDetailBack:
		general.RequireRole(app, session, HandleEntitiesCommand, roles.Admin)

	case states.CallbackAdminUserDetailAgain:
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)

	case states.CallbackAdminUserDetailLogs:
		general.RequireRole(app, session, handleUserLogs, roles.SuperAdmin)

	case states.CallbackAdminUserDetailUnban:
		general.RequireRole(app, session, handleUserUnban, roles.Admin)

	case states.CallbackAdminUserDetailBan:
		general.RequireRole(app, session, handleUserBan, roles.Admin)

	case states.CallbackAdminUserDetailRole:
		general.RequireRole(app, session, handleUserRole, roles.SuperAdmin)

	default:
		if strings.HasPrefix(callback, states.PrefixSelectAdminUserRole) {
			general.RequireRole(app, session, handleUserDetailSelect, roles.SuperAdmin)
		}
	}
}

func HandleUserDetailProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	switch session.State {
	case states.ProcessAdminUserDetailAwaitingReason:
		parseUserBanReason(app, session)
	}
}

func handleUserDetailSelect(app models.App, session *models.Session) {
	var role roles.Role

	switch utils.ParseCallback(app.Update) {
	case states.CallbackAdminUserRoleSelectUser:
		role = roles.User

	case states.CallbackAdminUserRoleSelectHelper:
		role = roles.Helper

	case states.CallbackAdminUserRoleSelectAdmin:
		role = roles.Admin

	case states.CallbackAdminUserRoleSelectSuper:
		role = roles.SuperAdmin
	}

	processUserRole(app, session, role)
}

func handleUserLogs(app models.App, session *models.Session) {
	if path, err := utils.GetLogFilePath(session.AdminState.UserID); err != nil {
		app.SendMessage(messages.LogsNotFound(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminUserDetailAgain))
	} else {
		app.SendFile(path, messages.LogsFound(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminUserDetailAgain))
	}
}

func handleUserUnban(app models.App, session *models.Session) {
	if err := postgres.SetUserBanStatus(session.AdminState.UserID, false); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminUserDetailAgain))
		return
	}

	app.SendMessage(messages.Unban(session), nil)
	app.SendMessageByID(session.AdminState.UserID, messages.UnbanNotification(session), nil)
	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func handleUserBan(app models.App, session *models.Session) {
	if session.AdminState.UserRole.HasAccess(roles.Helper) {
		app.SendMessage(messages.NeedRemoveRole(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminUserDetailAgain))
		return
	}

	app.SendMessage(messages.RequestBanReason(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessAdminUserDetailAwaitingReason)
}

func parseUserBanReason(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		processUserBan(app, session, "")
	} else {
		processUserBan(app, session, utils.ParseMessageString(app.Update))
	}
}

func processUserBan(app models.App, session *models.Session, reason string) {
	if err := postgres.SetUserBanStatus(session.AdminState.UserID, true); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminUserDetailAgain))
		return
	}

	app.SendMessage(messages.Ban(session, reason), nil)
	app.SendMessageByID(session.AdminState.UserID, messages.BanNotification(session, reason), nil)
	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func handleUserRole(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceRole(session), keyboards.BuildAdminUserRoleKeyboard(session))
}

func processUserRole(app models.App, session *models.Session, role roles.Role) {
	if !canChangeRole(session, role > session.AdminState.UserRole) {
		app.SendMessage(messages.NoAccess(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminUserDetailAgain))
		return
	}

	if err := postgres.SetUserRole(session.AdminState.UserID, role); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminUserDetailAgain))
		return
	}

	app.SendMessage(messages.ChangeRole(session, role), nil)
	app.SendMessageByID(session.AdminState.UserID, messages.ChangeRoleNotification(session, role), nil)
	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}
