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
	app.SendMessage(messages.BuildSortingMessage(session), keyboards.BuildFilmsSortingKeyboard(session))
}

func HandleSortingFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackSortingFilmsBack:
		handleWithResetFilmsPage(app, session, HandleFilmsCommand)
	case states.CallbackSortingFilmsAllReset:
		handleSortingFilmsReset(app, session, HandleFilmsCommand)
	default:
		handleSortingFilmsSelect(app, session, callback)
	}
}

func handleSortingFilmsSelect(app models.App, session *models.Session, callback string) {
	if strings.HasPrefix(callback, states.PrefixSortingFilmsSelect) {
		session.GetFilmsSortingByContext().Field = strings.TrimPrefix(callback, states.PrefixSortingFilmsSelect)
		handleSortingFilmsDirection(app, session)
	}
}

func HandleSortingFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		handleWithResetFilmsPage(app, session, HandleSortingFilmsCommand)
		return
	}

	switch session.State {
	case states.ProcessSortingFilmsAwaitingDirection:
		parseSortingFilmsDirection(app, session)
	}
}

func handleSortingFilmsDirection(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildSelectedSortMessage(session, session.GetFilmsSortingByContext()), keyboards.BuildSortingDirectionKeyboard(session, session.GetFilmsSortingByContext()))
	session.SetState(states.ProcessSortingFilmsAwaitingDirection)
}

func parseSortingFilmsDirection(app models.App, session *models.Session) {
	if utils.IsReset(app.Update) {
		handleSortingFilmsReset(app, session, HandleSortingFilmsCommand)
		return
	}

	if utils.IsDecrease(app.Update) {
		session.GetFilmsSortingByContext().Direction = "-"
	}

	session.GetFilmsSortingByContext().SetSort()
	app.SendMessage(messages.BuildSortingAppliedMessage(session, session.GetFilmsSortingByContext()), nil)
	handleWithResetFilmsPage(app, session, HandleFilmsCommand)
}

func handleSortingFilmsReset(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.GetFilmsSortingByContext().ResetSorting()
	app.SendMessage(messages.BuildSortingResetSuccessMessage(session), nil)
	handleWithResetFilmsPage(app, session, next)
}

func handleWithResetFilmsPage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	clearStatesAndResetFilmsPage(session)
	next(app, session)
}
