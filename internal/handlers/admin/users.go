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
	"log"
	"strconv"
	"strings"
)

func HandleUsersCommand(app models.App, session *models.Session) {
	//msg := translator.Translate(session.Lang, "choiceAction", nil, nil)
	//
	//keyboard := keyboards.BuildAdminUsersKeyboard(session)
	//
	//app.SendMessage(msg, keyboard)
	users, err := parseUsers(session)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleMenuCommand, roles.Admin)
		return
	}

	msg := messages.BuildAdminUserListMessage(session, users)

	keyboard := keyboards.BuildAdminUserListKeyboard(session, users)

	app.SendMessage(msg, keyboard)
}

func HandleUsersButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	switch {
	case callback == states.CallbackAdminManageUsersSelectBack:
		general.RequireRole(app, session, HandleMenuCommand, roles.Helper)

	case callback == states.CallbackAdminManageUsersSelectFind:
		general.RequireRole(app, session, handleUserFindCommand, roles.Admin)

	case callback == states.CallbackAdminUsersListNextPage:
		if session.AdminState.CurrentPage < session.AdminState.LastPage {
			session.AdminState.CurrentPage++
			HandleUsersCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminUsersListPrevPage:
		if session.AdminState.CurrentPage > 1 {
			session.AdminState.CurrentPage--
			HandleUsersCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case strings.HasPrefix(callback, "select_admin_user_"):
		handleUserSelect(app, session)
	}
}

func HandleUsersProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessAdminManageUsersAwaitingFind:
		general.RequireRole(app, session, processUserFindSelect, roles.Admin)
	}
}

func handleUserSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_admin_user_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		log.Printf("error parsing user ID: %v", err)
		return
	}

	session.AdminState.UserID = id

	HandleUserDetailCommand(app, session)
}

func handleUserFindCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "requestIDOrUsername", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessAdminManageUsersAwaitingFind)
}

func processUserFindSelect(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleUsersCommand(app, session)
		return
	}

	param := utils.ParseMessageString(app.Upd)

	if strings.HasPrefix(param, "@") {
		param = strings.TrimPrefix(param, "@")
		user, err := postgres.GetUserByTelegramUsername(param)
		if err != nil || user == nil {
			msg := translator.Translate(session.Lang, "notFound", nil, nil)
			app.SendMessage(msg, nil)
			handleUserFindCommand(app, session)
			return
		}
		session.AdminState.UserID = user.TelegramID
	} else {
		telegramID, err := strconv.Atoi(param)
		if err != nil {
			msg := translator.Translate(session.Lang, "someError", nil, nil)
			app.SendMessage(msg, nil)
			handleUserFindCommand(app, session)
			return
		}

		user, err := postgres.GetUserByTelegramID(telegramID)
		if err != nil || user == nil {
			msg := translator.Translate(session.Lang, "notFound", nil, nil)
			app.SendMessage(msg, nil)
			handleUserFindCommand(app, session)
			return
		}
		session.AdminState.UserID = user.TelegramID
	}

	session.ClearState()
	HandleUserDetailCommand(app, session)
}

func parseUsers(session *models.Session) ([]models.Session, error) {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize

	users, err := postgres.GetAllUsersWithPagination(currentPage, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := postgres.GetUserCounts()
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount / int64(pageSize))
	if totalCount%int64(pageSize) > 0 {
		totalPages++
	}

	session.AdminState.LastPage = totalPages
	session.AdminState.TotalRecords = int(totalCount)

	return users, nil
}
=======
//func HandleUsersCommand(app models.App, session *models.Session) {
//	msg := translator.Translate(session.Lang, "choiceAction", nil, nil)
//
//	//keyboard :=
//
//	// Список пользователей
//	// Поиск пользователя по айди или по нику
//
//}
//
//func HandleUsersButton(app models.App, session *models.Session) {
//
//}
//
//func HandleUserProcess(app models.App, session *models.Session) {
//
//}
