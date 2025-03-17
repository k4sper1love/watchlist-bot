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

// HandleAdminDetailCommand handles the command for viewing detailed information about an admin user.
// Retrieves the admin's details and sends a message with their information and an appropriate keyboard.
func HandleAdminDetailCommand(app models.App, session *models.Session) {
	if admin, err := getEntity(session); err != nil {
		resetAdminPageAndHandle(app, session, HandleEntitiesCommand, roles.Admin)
	} else {
		app.SendMessage(messages.UserDetail(session, admin), keyboards.AdminDetail(session, admin))
	}
}

// HandleAdminDetailButtons handles button interactions related to the admin detail view.
// Supports actions like going back, refreshing details, raising/lowering roles, and removing roles.
func HandleAdminDetailButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallAdminDetailBack:
		general.RequireRole(app, session, HandleEntitiesCommand, roles.Admin)

	case states.CallAdminDetailAgain:
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)

	case states.CallAdminDetailRaiseRole:
		handleRoleChange(app, session, true)

	case states.CallAdminDetailLowerRole:
		handleRoleChange(app, session, false)

	case states.CallAdminDetailRemoveRole:
		general.RequireRole(app, session, handleRemoveRole, roles.SuperAdmin)
	}
}

// handleRoleChange processes role changes (raising or lowering) for a user in the admin panel.
// Ensures the current user has sufficient permissions and updates the role in the database.
func handleRoleChange(app models.App, session *models.Session, raise bool) {
	if !canChangeRole(session, raise) {
		app.SendMessage(messages.NoAccess(session), keyboards.Back(session, states.CallAdminDetailAgain))
		return
	}

	if err := postgres.SetUserRole(session.AdminState.UserID, getNewRole(session.AdminState.UserRole, raise)); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallAdminDetailAgain))
		return
	}

	app.SendMessageByID(session.AdminState.UserID, messages.ShiftRoleNotification(session, raise), nil)
	app.SendMessage(messages.ShiftRole(session, raise), nil)
	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}

// canChangeRole checks if the current user is allowed to change the role of another user.
// Considers constraints like self-role changes, hierarchy rules, and role access levels.
func canChangeRole(session *models.Session, raise bool) bool {
	if session.AdminState.UserID == session.TelegramID {
		return false
	}
	if !session.Role.HasAccess(roles.SuperAdmin) {
		return false
	}
	if raise {
		return !(session.AdminState.UserRole.HasAccess(roles.Admin) && !session.Role.HasAccess(roles.Root) || session.AdminState.UserRole.HasAccess(roles.SuperAdmin))
	}
	return !(session.AdminState.UserRole.HasAccess(roles.SuperAdmin) && !session.Role.HasAccess(roles.Root) || session.AdminState.UserRole == roles.User)
}

// getNewRole calculates the new role based on the current role and whether the role is being raised or lowered.
func getNewRole(current roles.Role, raise bool) roles.Role {
	if raise {
		return current.NextRole()
	}
	return current.PrevRole()
}

// handleRemoveRole processes the removal of a user's role, reverting it to the default "User" role.
// Ensures the current user has sufficient permissions and updates the role in the database.
func handleRemoveRole(app models.App, session *models.Session) {
	if session.AdminState.UserRole.HasAccess(session.Role) && session.Role.HasAccess(roles.SuperAdmin) {
		app.SendMessage(messages.NoAccess(session), keyboards.Back(session, states.CallAdminDetailAgain))
		return
	}

	if err := postgres.SetUserRole(session.AdminState.UserID, roles.User); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallAdminDetailAgain))
		return
	}

	app.SendMessageByID(session.AdminState.UserID, messages.RemoveRoleNotification(session), nil)
	app.SendMessage(messages.RemoveRole(session), nil)
	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}
