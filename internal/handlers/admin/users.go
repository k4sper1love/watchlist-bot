package admin

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
	"strconv"
	"strings"
)

func HandleUsersCommand(app models.App, session *models.Session) {
	users, err := parseUsers(session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackMenuSelectAdmin).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := messages.BuildAdminUserListMessage(session, users)

	keyboard := keyboards.BuildAdminUserListKeyboard(session, users)

	app.SendMessage(msg, keyboard)
}

func HandleUsersButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)
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
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminUsersListPrevPage:
		if session.AdminState.CurrentPage > 1 {
			session.AdminState.CurrentPage--
			HandleUsersCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminUsersListLastPage:
		if session.AdminState.CurrentPage != session.AdminState.LastPage {
			session.AdminState.CurrentPage = session.AdminState.LastPage
			HandleUsersCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminUsersListFirstPage:
		if session.AdminState.CurrentPage != 1 {
			session.AdminState.CurrentPage = 1
			HandleUsersCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case strings.HasPrefix(callback, "select_admin_user_"):
		handleUserSelect(app, session)
	}
}

func HandleUsersProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleUsersCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAdminManageUsersAwaitingFind:
		general.RequireRole(app, session, processUserFindSelect, roles.Admin)
	}
}

func handleUserSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)
	idStr := strings.TrimPrefix(callback, "select_admin_user_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminSelectUsers).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		sl.Log.Error("failed to parse user ID", slog.Any("error", err), slog.String("callback", callback))
		return
	}

	session.AdminState.UserID = id

	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func handleUserFindCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "requestIDOrUsername", nil, nil)

	hintMsg := translator.Translate(session.Lang, "hintAPIUserID", nil, nil)
	msg += fmt.Sprintf("\n\n<i>%s</i>", hintMsg)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessAdminManageUsersAwaitingFind)
}

func processUserFindSelect(app models.App, session *models.Session) {
	param := utils.ParseMessageString(app.Update)

	if strings.HasPrefix(param, "@") {
		param = strings.TrimPrefix(param, "@")
		user, err := postgres.GetUserByField(postgres.TelegramUsernameField, param, false)
		if err != nil || user == nil {
			msg := "❗" + translator.Translate(session.Lang, "notFound", nil, nil)
			app.SendMessage(msg, nil)
			handleUserFindCommand(app, session)
			return
		}
		session.AdminState.UserID = user.TelegramID

	} else if strings.HasPrefix(param, "api_") {
		param = strings.TrimPrefix(param, "api_")
		id, err := handleAndParseID(app, session, param)
		if err != nil {
			handleUserFindCommand(app, session)
			return
		}

		user, err := postgres.GetUserByAPIUserID(id)
		if err != nil {
			msg := "❗" + translator.Translate(session.Lang, "notFound", nil, nil)
			app.SendMessage(msg, nil)
			handleUserFindCommand(app, session)
			return
		}

		session.AdminState.UserID = user.TelegramID

	} else {
		telegramID, err := handleAndParseID(app, session, param)
		if err != nil {
			handleUserFindCommand(app, session)
			return
		}

		user, err := postgres.GetUserByField(postgres.TelegramIDField, telegramID, false)
		if err != nil || user == nil {
			msg := "❗️" + translator.Translate(session.Lang, "notFound", nil, nil)
			app.SendMessage(msg, nil)
			handleUserFindCommand(app, session)
			return
		}
		session.AdminState.UserID = user.TelegramID
	}

	session.ClearState()
	general.RequireRole(app, session, HandleUserDetailCommand, roles.Admin)
}

func parseUsers(session *models.Session) ([]models.Session, error) {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize

	users, err := postgres.GetUsersWithPagination(currentPage, pageSize, false)
	if err != nil {
		return nil, err
	}

	totalCount, err := postgres.GetUserCount(false)
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

func handleAndParseID(app models.App, session *models.Session, param string) (int, error) {
	parsed, err := strconv.Atoi(param)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		return -1, err
	}

	return parsed, nil
}
