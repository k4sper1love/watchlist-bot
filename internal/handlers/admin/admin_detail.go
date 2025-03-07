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
		app.SendMessage(messages.BuildSomeErrorMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminSelectAdmins))
	} else {
		app.SendMessage(messages.BuildAdminUserDetailMessage(session, admin), keyboards.BuildAdminDetailKeyboard(session, admin))
	}
}

func HandleAdminDetailButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackAdminDetailBack:
		general.RequireRole(app, session, HandleEntitiesCommand, roles.Admin)

	case states.CallbackAdminDetailAgain:
		general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)

	case states.CallbackAdminDetailRaiseRole:
		handleRoleChange(app, session, true)

	case states.CallbackAdminDetailLowerRole:
		handleRoleChange(app, session, false)

	case states.CallbackAdminDetailRemoveRole:
		general.RequireRole(app, session, handleRemoveRole, roles.SuperAdmin)
	}
}

func handleRoleChange(app models.App, session *models.Session, raise bool) {
	if !canChangeRole(session, raise) {
		app.SendMessage(messages.BuildNoAccessMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminDetailAgain))
		return
	}

	if err := postgres.SetUserRole(session.AdminState.UserID, getNewRole(session.AdminState.UserRole, raise)); err != nil {
		app.SendMessage(messages.BuildSomeErrorMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminDetailAgain))
		return
	}

	notifyUserAboutRole(app, session, raise)
	app.SendMessage(messages.BuildRoleChangeMessage(session, raise), nil)
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
	return !(session.AdminState.UserRole.HasAccess(roles.SuperAdmin) && !session.Role.HasAccess(roles.Root) || session.AdminState.UserRole.HasAccess(roles.User))
}

func getNewRole(current roles.Role, raise bool) roles.Role {
	if raise {
		return current.NextRole()
	}
	return current.PrevRole()
}

func notifyUserAboutRole(app models.App, session *models.Session, raise bool) {
	if raise {
		app.SendMessageByID(session.AdminState.UserID, messages.BuildRaiseRoleNotificationMessage(session), nil)
	} else {
		app.SendMessageByID(session.AdminState.UserID, messages.BuildLowerRoleNotificationMessage(session), nil)
	}
}

func handleRemoveRole(app models.App, session *models.Session) {
	if session.AdminState.UserRole.HasAccess(session.Role) {
		app.SendMessage(messages.BuildNoAccessMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminDetailAgain))
		return
	}

	if err := postgres.SetUserRole(session.AdminState.UserID, roles.User); err != nil {
		app.SendMessage(messages.BuildSomeErrorMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackAdminDetailAgain))
		return
	}

	app.SendMessageByID(session.AdminState.UserID, messages.BuildRemoveRoleNotificationMessage(session), nil)
	app.SendMessage(messages.BuildRemoveRoleMessage(session), nil)
	general.RequireRole(app, session, HandleAdminDetailCommand, roles.Admin)
}
