package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleSortingFilmsCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "choiceSorting", nil, nil)

	keyboard := keyboards.BuildFilmsSortingKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleSortingFilmsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackSortingFilmsSelectBack:
		HandleFilmsCommand(app, session)
		return

	case states.CallbackSortingFilmsSelectAllReset:
		handleSortingFilmsAllReset(app, session)
		return

	case states.CallbackSortingFilmsSelectID:
		session.GetFilmsSortingByContext().Field = "id"

	case states.CallbackSortingFilmsSelectTitle:
		session.GetFilmsSortingByContext().Field = "title"

	case states.CallbackSortingFilmsSelectRating:
		session.GetFilmsSortingByContext().Field = "rating"
	}

	handleSortingFilmsDirection(app, session)
}

func HandleSortingFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleSortingFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessSortingFilmsAwaitingDirection:
		parseSortingFilmsDirection(app, session)
	}

}

func handleSortingFilmsAllReset(app models.App, session *models.Session) {
	session.GetFilmsSortingByContext().ResetSorting()

	msg := translator.Translate(session.Lang, "sortingResetSuccess", nil, nil)

	app.SendMessage(msg, nil)

	HandleSortingFilmsCommand(app, session)
}

func handleSortingFilmsDirection(app models.App, session *models.Session) {
	msg := messages.BuildSelectedSortMessage(session)

	keyboard := keyboards.BuildFilmsSortingDirectionKeyboard(session)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSortingFilmsAwaitingDirection)
}

func parseSortingFilmsDirection(app models.App, session *models.Session) {
	sorting := session.GetFilmsSortingByContext()

	if utils.IsReset(app.Upd) {
		sorting.Sort = ""
		handleSortingFilmsReset(app, session)
		return
	}

	if utils.ParseCallback(app.Upd) == states.CallbacktDecrease {
		sorting.Direction = "-"
	}
	sorting.Sort = sorting.Direction + sorting.Field

	msg := translator.Translate(session.Lang, "sortingApplied", nil, nil)
	app.SendMessage(msg, nil)

	session.ClearAllStates()

	HandleSortingFilmsCommand(app, session)
}

func handleSortingFilmsReset(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "sortingResetSuccess", nil, nil)
	app.SendMessage(msg, nil)

	session.ClearAllStates()
	HandleSortingFilmsCommand(app, session)
}
