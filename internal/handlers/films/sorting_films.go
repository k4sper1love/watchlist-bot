package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

func HandleSortingFilmsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceSorting(session), keyboards.FilmsSorting(session))
}

func HandleSortingFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFilmSortingBack:
		handleWithResetFilmsPage(app, session, HandleFilmsCommand)
	case states.CallFilmSortingAllReset:
		handleSortingFilmsReset(app, session, HandleFilmsCommand)
	default:
		handleSortingFilmsSelect(app, session, callback)
	}
}

func handleSortingFilmsSelect(app models.App, session *models.Session, callback string) {
	if strings.HasPrefix(callback, states.FilmSortingSelect) {
		session.GetFilmSortingByCtx().Field = strings.TrimPrefix(callback, states.FilmSortingSelect)
		handleSortingFilmsDirection(app, session)
	}
}

func HandleSortingFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		handleWithResetFilmsPage(app, session, HandleSortingFilmsCommand)
		return
	}

	switch session.State {
	case states.AwaitFilmSortingDirection:
		parseSortingFilmsDirection(app, session)
	}
}

func handleSortingFilmsDirection(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestSortDirection(session, session.GetFilmSortingByCtx()), keyboards.SortingDirection(session, session.GetFilmSortingByCtx()))
	session.SetState(states.AwaitFilmSortingDirection)
}

func parseSortingFilmsDirection(app models.App, session *models.Session) {
	if utils.IsReset(app.Update) {
		handleSortingFilmsReset(app, session, HandleSortingFilmsCommand)
		return
	}

	if utils.IsDecrease(app.Update) {
		session.GetFilmSortingByCtx().Direction = "-"
	}

	session.GetFilmSortingByCtx().SetSort()
	app.SendMessage(messages.SortingApplied(session, session.GetFilmSortingByCtx()), nil)
	handleWithResetFilmsPage(app, session, HandleFilmsCommand)
}

func handleSortingFilmsReset(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.GetFilmSortingByCtx().Reset()
	app.SendMessage(messages.ResetSortingSuccess(session), nil)
	handleWithResetFilmsPage(app, session, next)
}

func handleWithResetFilmsPage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	clearStatesAndResetFilmsPage(session)
	next(app, session)
}
