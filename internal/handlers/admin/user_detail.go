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
	user, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleUsersCommand(app, session)
		return
	}

	msg := messages.BuildAdminUserDetailMessage(session, user)
	keyboard := keyboards.BuildAdminUserDetailKeyboard(session, user)

	app.SendMessage(msg, keyboard)
}

func HandleUserDetailButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackAdminUserDetailBack:
		general.RequireRole(app, session, HandleUsersCommand, roles.Admin)

	case callback == states.CallbackAdminUserDetailUnban:
		general.RequireRole(app, session, handleUserUnban, roles.Admin)

	case callback == states.CallbackAdminUserDetailBan:
		general.RequireRole(app, session, handleUserBan, roles.Admin)

	case callback == states.CallbackAdminUserDetailRole:
		general.RequireRole(app, session, handleUserRole, roles.SuperAdmin)

	//case callback == states.CallbackAdminUserDetailFeedback:

	case strings.HasPrefix(callback, "admin_user_role_select_"):
		general.RequireRole(app, session, processUserRole, roles.SuperAdmin)
	}
}

func handleUserUnban(app models.App, session *models.Session) {
	err := postgres.UnbanUser(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
	} else {
		msg := translator.Translate(session.Lang, "unbanned", nil, nil)
		app.SendMessage(msg, nil)
	}

	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func handleUserBan(app models.App, session *models.Session) {
	user, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	if user.Role.HasAccess(roles.Helper) {
		msg := translator.Translate(session.Lang, "needRemoveRole", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	err = postgres.BanUser(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
	} else {
		msg := translator.Translate(session.Lang, "banned", nil, nil)
		app.SendMessage(msg, nil)
	}

	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func handleUserRole(app models.App, session *models.Session) {
	user, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	part1 := translator.Translate(session.Lang, "currentRole", nil, nil)
	part2 := translator.Translate(session.Lang, "choiceRole", nil, nil)
	msg := fmt.Sprintf("%s: %s\n\n%s", part1, user.Role.String(), part2)

	keyboard := keyboards.BuildAdminUserRoleKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func processUserRole(app models.App, session *models.Session) {
	user, err := postgres.GetUserByTelegramID(session.AdminState.UserID)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	if user.Role == roles.SuperAdmin && !session.Role.HasAccess(roles.Root) {
		msg := translator.Translate(session.Lang, "noAccess", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	if user.Role.HasAccess(roles.Root) {
		msg := translator.Translate(session.Lang, "noAccess", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	var role roles.Role

	switch utils.ParseCallback(app.Upd) {
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

	_, err = postgres.SetUserRole(user.TelegramID, role)
	if err != nil {
		msg = translator.Translate(session.Lang, "someError", nil, nil)
	} else {
		msg = translator.Translate(session.Lang, "success", nil, nil)
	}

	app.SendMessage(msg, nil)

	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}
