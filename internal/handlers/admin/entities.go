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
	"strconv"
	"strings"
)

// HandleEntitiesCommand handles the command for listing entities (users or admins).
// Retrieves paginated entities and sends a message with their details and navigation keyboard.
func HandleEntitiesCommand(app models.App, session *models.Session) {
	if entities, err := getEntities(session); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallMenuAdmin))
	} else {
		app.SendMessage(messages.UserList(session, entities), keyboards.AdminList(session, entities))
	}
}

// HandleEntitiesButtons handles button interactions related to entity management.
// Supports actions like going back, searching, pagination, and selecting specific entities.
func HandleEntitiesButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallEntitiesBack:
		general.RequireRole(app, session, HandleMenuCommand, roles.Helper)

	case states.CallEntitiesFind:
		general.RequireRole(app, session, handleEntitiesFindCommand, roles.Admin)

	default:
		if strings.HasPrefix(callback, states.EntitiesPage) {
			handleEntitiesPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectEntity) {
			handleEntitiesSelect(app, session, callback)
		}
	}
}

// HandleEntitiesProcess processes the entity-related workflow based on the current session state.
// Handles states like awaiting a search query for entities.
func HandleEntitiesProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.AwaitEntitiesFind:
		general.RequireRole(app, session, parseEntitiesFind, roles.Admin)
	}
}

// handleEntitiesPagination processes pagination actions for entity lists.
// Updates the current page in the session and reloads the entity list.
func handleEntitiesPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallEntitiesPageNext:
		if session.AdminState.CurrentPage >= session.AdminState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage++

	case states.CallEntitiesPagePrev:
		if session.AdminState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage--

	case states.CallEntitiesPageLast:
		if session.AdminState.CurrentPage == session.AdminState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage = session.AdminState.LastPage

	case states.CallEntitiesFirst:
		if session.AdminState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage = 1
	}

	HandleEntitiesCommand(app, session)
}

// handleEntitiesSelect processes the selection of an entity from the list.
// Parses the entity ID and navigates to the entity detail view.
func handleEntitiesSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, getSelectEntityPrefix(session))); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallMenuAdmin))
	} else {
		session.AdminState.UserID = id
		HandleEntityDetailCommand(app, session)
	}
}

// HandleEntityDetailCommand handles the command for viewing detailed information about a selected entity.
// Delegates to either admin or user detail commands based on the entity type.
func HandleEntityDetailCommand(app models.App, session *models.Session) {
	if session.AdminState.IsAdmin {
		HandleAdminDetailCommand(app, session)
	} else {
		HandleUserDetailCommand(app, session)
	}
}

// handleEntitiesFindCommand prompts the user to enter a search query for finding entities.
func handleEntitiesFindCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestEntityField(session), keyboards.Cancel(session))
	session.SetState(states.AwaitEntitiesFind)
}

// parseEntitiesFind processes the search query for entities.
// Retrieves the entity by Telegram ID, username, or API user ID and navigates to its detail view.
func parseEntitiesFind(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleEntitiesCommand(app, session)
		return
	}

	user, err := parseEntityByField(session, utils.ParseMessageString(app.Update))
	if err != nil || user == nil {
		app.SendMessage(messages.NotFound(session), nil)
		handleEntitiesFindCommand(app, session)
		return
	}

	session.AdminState.UserID = user.TelegramID
	session.ClearState()
	HandleEntityDetailCommand(app, session)
}

// parseEntityByField retrieves an entity based on the provided field (username, Telegram ID, or API user ID).
func parseEntityByField(session *models.Session, input string) (*models.Session, error) {
	switch {
	case strings.HasPrefix(input, "@"):
		return postgres.GetUserByField(postgres.TelegramUsernameField, strings.TrimPrefix(input, "@"), session.AdminState.IsAdmin)

	case strings.HasPrefix(input, "api_"):
		id, err := strconv.Atoi(strings.TrimPrefix(input, "api_"))
		if err != nil {
			return nil, err
		}
		return postgres.GetUserByAPIUserID(id, session.AdminState.IsAdmin)

	default:
		telegramID, err := strconv.Atoi(input)
		if err != nil {
			return nil, err
		}
		return postgres.GetUserByField(postgres.TelegramIDField, telegramID, session.AdminState.IsAdmin)
	}
}

// getSelectEntityPrefix determines the prefix for selecting entities based on the entity type (admin or user).
func getSelectEntityPrefix(session *models.Session) string {
	if session.AdminState.IsAdmin {
		return states.SelectAdmin
	}
	return states.SelectUser
}

// getEntities retrieves paginated entities (users or admins) from the database.
// Calculates pagination metadata and returns the entities.
func getEntities(session *models.Session) ([]models.Session, error) {
	users, err := postgres.GetUsersWithPagination(session.AdminState.CurrentPage, session.AdminState.PageSize, session.AdminState.IsAdmin)
	if err != nil {
		return nil, err
	}

	totalCount, err := postgres.GetUserCount(session.AdminState.IsAdmin)
	if err != nil {
		return nil, err
	}

	calculateAdminPages(session, totalCount)
	return users, nil
}

// getEntity retrieves a specific entity by its Telegram ID.
// Updates the session with the entity's details.
func getEntity(session *models.Session) (*models.Session, error) {
	entity, err := postgres.GetUserByField(postgres.TelegramIDField, session.AdminState.UserID, session.AdminState.IsAdmin)
	if err != nil {
		return nil, err
	}

	updateSessionWithEntity(session, entity)
	return entity, nil
}

// updateSessionWithEntity updates the session with the selected entity's details (language and role).
func updateSessionWithEntity(session *models.Session, entity *models.Session) {
	session.AdminState.UserLang = entity.Lang
	session.AdminState.UserRole = entity.Role
}

// calculateAdminPages calculates pagination metadata (last page and total records) for entity lists.
func calculateAdminPages(session *models.Session, total int64) {
	session.AdminState.LastPage = (int(total) + session.AdminState.PageSize - 1) / session.AdminState.PageSize
	session.AdminState.TotalRecords = int(total)
}
