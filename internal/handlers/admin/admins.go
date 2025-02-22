package admin

import (
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

func HandleAdminsCommand(app models.App, session *models.Session) {
	admins, err := parseAdmins(session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackMenuSelectAdmin).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := messages.BuildAdminUserListMessage(session, admins)

	keyboard := keyboards.BuildAdminListKeyboard(session, admins)

	app.SendMessage(msg, keyboard)
}

func HandleAdminsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	switch {
	case callback == states.CallbackAdminListBack:
		general.RequireRole(app, session, HandleMenuCommand, roles.Helper)

	case callback == states.CallbackAdminListSelectFind:
		general.RequireRole(app, session, handleAdminFindCommand, roles.Admin)

	case callback == states.CallbackAdminListNextPage:
		if session.AdminState.CurrentPage < session.AdminState.LastPage {
			session.AdminState.CurrentPage++
			HandleAdminsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminListPrevPage:
		if session.AdminState.CurrentPage > 1 {
			session.AdminState.CurrentPage--
			HandleAdminsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminListLastPage:
		if session.AdminState.CurrentPage != session.AdminState.LastPage {
			session.AdminState.CurrentPage = session.AdminState.LastPage
			HandleAdminsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminListFirstPage:
		if session.AdminState.CurrentPage != 1 {
			session.AdminState.CurrentPage = 1
			HandleAdminsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case strings.HasPrefix(callback, "select_admin_"):
		handleAdminsSelect(app, session)
	}
}

func HandleAdminsProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessAdminListAwaitingFind:
		general.RequireRole(app, session, processAdminFindSelect, roles.Admin)
	}
}

func handleAdminsSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_admin_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackMenuSelectAdmin).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		sl.Log.Warn("failed to parse admin ID", slog.Any("error", err), slog.String("callback", callback))
		return
	}

	session.AdminState.UserID = id

	HandleAdminDetailCommand(app, session)
}

func handleAdminFindCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "requestIDOrUsername", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessAdminListAwaitingFind)
}

func processAdminFindSelect(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleAdminsCommand(app, session)
		return
	}

	param := utils.ParseMessageString(app.Upd)

	if strings.HasPrefix(param, "@") {
		param = strings.TrimPrefix(param, "@")
		user, err := postgres.GetUserByField(postgres.TelegramUsernameField, param, true)
		if err != nil || user == nil {
			msg := "❗️" + translator.Translate(session.Lang, "notFound", nil, nil)
			app.SendMessage(msg, nil)
			handleAdminFindCommand(app, session)
			return
		}
		session.AdminState.UserID = user.TelegramID
	} else {
		telegramID, err := strconv.Atoi(param)
		if err != nil {
			msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
			keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminListSelectFind).Build(session.Lang)
			app.SendMessage(msg, keyboard)
			return
		}

		user, err := postgres.GetUserByField(postgres.TelegramIDField, telegramID, true)
		if err != nil || user == nil {
			msg := "❗️" + translator.Translate(session.Lang, "notFound", nil, nil)
			app.SendMessage(msg, nil)
			handleAdminFindCommand(app, session)
			return
		}
		session.AdminState.UserID = user.TelegramID
	}

	session.ClearState()
	HandleAdminDetailCommand(app, session)
}

func parseAdmins(session *models.Session) ([]models.Session, error) {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize

	admins, err := postgres.GetUsersWithPagination(currentPage, pageSize, true)
	if err != nil {
		return nil, err
	}

	totalCount, err := postgres.GetUserCount(true)
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount / int64(pageSize))
	if totalCount%int64(pageSize) > 0 {
		totalPages++
	}

	session.AdminState.LastPage = totalPages
	session.AdminState.TotalRecords = int(totalCount)

	return admins, nil
}
