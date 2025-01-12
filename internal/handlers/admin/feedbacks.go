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

func HandleFeedbacksCommand(app models.App, session *models.Session) {
	feedbacks, err := parseFeedbacks(session)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		general.RequireRole(app, session, HandleMenuCommand, roles.Admin)
		return
	}

	msg := messages.BuildFeedbackListMessage(session, feedbacks)

	keyboard := keyboards.BuildFeedbackListKeyboard(session, feedbacks)

	app.SendMessage(msg, keyboard)
}

func HandleFeedbacksButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	switch {
	case callback == states.CallbackAdminFeedbackListBack:
		general.RequireRole(app, session, HandleMenuCommand, roles.Helper)

	case callback == states.CallbackAdminFeedbackListNextPage:
		if session.AdminState.CurrentPage < session.AdminState.LastPage {
			session.AdminState.CurrentPage++
			HandleFeedbacksCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminFeedbackListPrevPage:
		if session.AdminState.CurrentPage > 1 {
			session.AdminState.CurrentPage--
			HandleFeedbacksCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminFeedbackListLastPage:
		if session.AdminState.CurrentPage != session.AdminState.LastPage {
			session.AdminState.CurrentPage = session.AdminState.LastPage
			HandleFeedbacksCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAdminFeedbackListFirstPage:
		if session.AdminState.CurrentPage != 1 {
			session.AdminState.CurrentPage = 1
			HandleFeedbacksCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case strings.HasPrefix(callback, "select_admin_feedback_"):
		handleFeedbackSelect(app, session)
	}

}

func HandleFeedbacksProcess(app models.App, session *models.Session) {

}

func handleFeedbackSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_admin_feedback_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		log.Printf("error parsing user ID: %v", err)
		return
	}

	session.AdminState.FeedbackID = id

	HandleFeedbackDetailCommand(app, session)
}

func parseFeedbacks(session *models.Session) ([]models.Feedback, error) {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize

	feedback, err := postgres.GetAllFeedbacksWithPagination(currentPage, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := postgres.GetFeedbackCounts()
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount / int64(pageSize))
	if totalCount%int64(pageSize) > 0 {
		totalPages++
	}

	session.AdminState.LastPage = totalPages
	session.AdminState.TotalRecords = int(totalCount)

	return feedback, nil
}
