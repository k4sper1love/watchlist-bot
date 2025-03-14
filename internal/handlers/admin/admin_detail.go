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

func HandleAdminDetailCommand(app models.App, session *models.Session) {
	if admin, err := getEntity(session); err != nil {
		resetAdminPageAndHandle(app, session, HandleEntitiesCommand, roles.Admin)
	} else {
		app.SendMessage(messages.UserDetail(session, admin), keyboards.AdminDetail(session, admin))
	}
}

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

func getNewRole(current roles.Role, raise bool) roles.Role {
	if raise {
		return current.NextRole()
	}
	return current.PrevRole()
}

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
