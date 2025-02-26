package admin

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strings"
)

func HandleUserDetailCommand(app models.App, session *models.Session) {
	if user, err := getEntity(session); err != nil {
		app.SendMessage(messages.BuildSomeErrorMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminSelectUsers))
	} else {
		app.SendMessage(messages.BuildAdminUserDetailMessage(session, user), keyboards.BuildAdminUserDetailKeyboard(session, user))
	}
}

func HandleUserDetailButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackAdminUserDetail:
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)

	case states.CallbackAdminUserDetailBack:
		general.RequireRole(app, session, HandleEntitiesCommand, roles.Admin)

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
			general.RequireRole(app, session, processUserRole, roles.SuperAdmin)
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
		processUserBan(app, session)
	}
}

func handleUserLogs(app models.App, session *models.Session) {
	path, err := utils.GetLogFilePath(session.AdminState.UserID)
	if err != nil {
		msg := "❗" + translator.Translate(session.Lang, "logsNotFound", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	msg := "💾 " + translator.Translate(session.Lang, "logsFound", map[string]interface{}{
		"ID": session.AdminState.UserID,
	}, nil)

	keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminUserDetail).Build(session.Lang)

	app.SendFile(path, msg, keyboard)
}

func handleUserUnban(app models.App, session *models.Session) {
	err := postgres.SetUserBanStatus(session.AdminState.UserID, false)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
	} else {
		msg := messages.BuildUnbanMessage(session)
		app.SendMessage(msg, nil)

		msg = messages.BuildUserUnbanNotificationMessage(session)
		app.SendMessageByID(session.AdminState.UserID, msg, nil)
	}

	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func handleUserBan(app models.App, session *models.Session) {
	if session.AdminState.UserRole.HasAccess(roles.Helper) {
		msg := "❗" + translator.Translate(session.Lang, "needRemoveRole", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	msg := "❓" + translator.Translate(session.Lang, "requestBanReason", nil, nil)
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAdminUserDetailAwaitingReason)
}

func processUserBan(app models.App, session *models.Session) {
	var reason string

	if !utils.IsSkip(app.Update) {
		reason = utils.ParseMessageString(app.Update)
	}

	err := postgres.SetUserBanStatus(session.AdminState.UserID, true)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
	} else {
		msg := messages.BuildBanMessage(session, reason)
		app.SendMessage(msg, nil)

		msg = messages.BuildUserBanNotificationMessage(session, reason)
		app.SendMessageByID(session.AdminState.UserID, msg, nil)
	}

	session.ClearAllStates()

	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func handleUserRole(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "currentRole", nil, nil)
	part2 := translator.Translate(session.Lang, session.AdminState.UserRole.String(), nil, nil)
	part3 := translator.Translate(session.Lang, "choiceRole", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>: %s\n\n%s", part1, part2, part3)

	keyboard := keyboards.BuildAdminUserRoleKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func processUserRole(app models.App, session *models.Session) {
	if (session.AdminState.UserRole == roles.SuperAdmin && !session.Role.HasAccess(roles.Root)) || session.AdminState.UserRole.HasAccess(roles.Root) {
		msg := "❗" + translator.Translate(session.Lang, "noAccess", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	var role roles.Role

	switch utils.ParseCallback(app.Update) {
	case states.CallbackAdminUserRoleSelectBack:
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return

	case states.CallbackAdminUserRoleSelectUser:
		role = roles.User

	case states.CallbackAdminUserRoleSelectHelper:
		role = roles.Helper

	case states.CallbackAdminUserRoleSelectAdmin:
		role = roles.Admin

	case states.CallbackAdminUserRoleSelectSuper:
		role = roles.SuperAdmin
	}

	msg := " "

	err := postgres.SetUserRole(session.AdminState.UserID, role)
	if err != nil {
		msg = "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
	} else {
		msg = messages.BuildChangeRoleNotificationMessage(session, role)
		app.SendMessageByID(session.AdminState.UserID, msg, nil)

		msg = messages.BuildChangeRoleMessage(session, role)
	}

	app.SendMessage(msg, nil)

	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}
