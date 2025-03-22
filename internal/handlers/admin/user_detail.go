package admin

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/logger"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"strings"
)

// HandleUserDetailCommand handles the command for viewing detailed information about a specific user.
// Retrieves the user's details and sends a message with their information and an appropriate keyboard.
func HandleUserDetailCommand(app models.App, session *models.Session) {
	if user, err := getEntity(session); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallAdminUsers))
	} else {
		app.SendMessage(messages.UserDetail(session, user), keyboards.UserDetail(session, user))
	}
}

// HandleUserDetailButton handles button interactions related to the user detail view.
// Supports actions like going back, refreshing details, viewing logs, banning/unbanning, and changing roles.
func HandleUserDetailButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallUserDetailBack:
		general.RequireRole(app, session, HandleEntitiesCommand, roles.Admin)

	case states.CallUserDetailAgain:
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)

	case states.CallUserDetailLogs:
		general.RequireRole(app, session, handleUserLogs, roles.SuperAdmin)

	case states.CallUserDetailUnban:
		general.RequireRole(app, session, handleUserUnban, roles.Admin)

	case states.CallUserDetailBan:
		general.RequireRole(app, session, handleUserBan, roles.Admin)

	case states.CallUserDetailRole:
		general.RequireRole(app, session, handleUserRole, roles.SuperAdmin)

	default:
		if strings.HasPrefix(callback, states.CallUserDetailRole) {
			general.RequireRole(app, session, handleUserDetailSelect, roles.SuperAdmin)
		}
	}
}

// HandleUserDetailProcess processes the user detail workflow based on the current session state.
// Handles states like awaiting a ban reason.
func HandleUserDetailProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
		return
	}

	switch session.State {
	case states.AwaitUserDetailReason:
		parseUserBanReason(app, session)
	}
}

// handleUserDetailSelect processes the selection of a new role for the user.
func handleUserDetailSelect(app models.App, session *models.Session) {
	var role roles.Role

	switch utils.ParseCallback(app.Update) {
	case states.CallUserDetailRoleSelectUser:
		role = roles.User

	case states.CallUserDetailRoleSelectHelper:
		role = roles.Helper

	case states.CallUserDetailRoleSelectAdmin:
		role = roles.Admin

	case states.CallUserDetailRoleSelectSuper:
		role = roles.SuperAdmin
	}

	processUserRole(app, session, role)
}

// handleUserLogs retrieves and sends the log file for the specified user.
func handleUserLogs(app models.App, session *models.Session) {
	if path, err := logger.GetFilePath(session.AdminState.UserID); err != nil {
		app.SendMessage(messages.LogsNotFound(session), keyboards.Back(session, states.CallUserDetailAgain))
	} else {
		app.SendFile(path, messages.LogsFound(session), keyboards.Back(session, states.CallUserDetailAgain))
	}
}

// handleUserUnban unblocks the user by updating their ban status in the database.
func handleUserUnban(app models.App, session *models.Session) {
	if err := postgres.SetUserBanStatus(session.AdminState.UserID, false); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallUserDetailAgain))
		return
	}

	app.SendMessage(messages.Unban(session), nil)
	app.SendMessageByID(session.AdminState.UserID, messages.UnbanNotification(session), nil)
	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

// handleUserBan prompts the admin to provide a reason for banning the user.
func handleUserBan(app models.App, session *models.Session) {
	if session.AdminState.UserRole.HasAccess(roles.Helper) {
		app.SendMessage(messages.NeedRemoveRole(session), keyboards.Back(session, states.CallUserDetailAgain))
		return
	}

	app.SendMessage(messages.RequestBanReason(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitUserDetailReason)
}

// parseUserBanReason processes the ban reason provided by the admin.
// If skipped, bans the user without a reason; otherwise, uses the provided reason.
func parseUserBanReason(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		processUserBan(app, session, "")
	} else {
		processUserBan(app, session, utils.ParseMessageString(app.Update))
	}
}

// processUserBan updates the user's ban status in the database and sends notifications.
func processUserBan(app models.App, session *models.Session, reason string) {
	if err := postgres.SetUserBanStatus(session.AdminState.UserID, true); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallUserDetailAgain))
		return
	}

	app.SendMessage(messages.Ban(session, reason), nil)
	app.SendMessageByID(session.AdminState.UserID, messages.BanNotification(session, reason), nil)
	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

// handleUserRole displays the role selection menu for updating the user's role.
func handleUserRole(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceRole(session), keyboards.UserRoleSelect(session))
}

// processUserRole updates the user's role in the database and sends notifications.
func processUserRole(app models.App, session *models.Session, role roles.Role) {
	if !canChangeRole(session, role > session.AdminState.UserRole) {
		app.SendMessage(messages.NoAccess(session), keyboards.Back(session, states.CallUserDetailAgain))
		return
	}

	if err := postgres.SetUserRole(session.AdminState.UserID, role); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallUserDetailAgain))
		return
	}

	app.SendMessage(messages.ChangeRole(session, role), nil)
	app.SendMessageByID(session.AdminState.UserID, messages.ChangeRoleNotification(session, role), nil)
	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}
