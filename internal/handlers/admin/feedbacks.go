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

func HandleFeedbacksCommand(app models.App, session *models.Session) {
	if feedbacks, err := getFeedbacks(session); err != nil {
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallMenuAdmin))
	} else {
		app.SendMessage(messages.FeedbackList(session, feedbacks), keyboards.FeedbackList(session, feedbacks))
	}
}

func HandleFeedbacksButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFeedbacksBack:
		general.RequireRole(app, session, HandleMenuCommand, roles.Helper)

	default:
		if strings.HasPrefix(callback, states.FeedbacksPage) {
			handleFeedbackPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectFeedback) {
			handleFeedbackSelect(app, session, callback)
		}
	}
}

func handleFeedbackPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallFeedbacksPageNext:
		if session.AdminState.CurrentPage >= session.AdminState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage++

	case states.CallFeedbacksPagePrev:
		if session.AdminState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage--

	case states.CallFeedbacksPageLast:
		if session.AdminState.CurrentPage == session.AdminState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage = session.AdminState.LastPage

	case states.CallFeedbacksPageFirst:
		if session.AdminState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.AdminState.CurrentPage = 1
	}

	HandleFeedbacksCommand(app, session)
}

func handleFeedbackSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectFeedback)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.SomeError(session), keyboards.Back(session, states.CallAdminFeedback))
	} else {
		session.AdminState.FeedbackID = id
		HandleFeedbackDetailCommand(app, session)
	}
}

func getFeedbacks(session *models.Session) ([]models.Feedback, error) {
	feedback, err := postgres.GetFeedbacksWithPagination(session.AdminState.CurrentPage, session.AdminState.PageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := postgres.GetFeedbackCount()
	if err != nil {
		return nil, err
	}

	calculateAdminPages(session, totalCount)
	return feedback, nil
}
