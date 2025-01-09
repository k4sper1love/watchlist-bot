package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleSortingCollectionsCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "choiceSorting", nil, nil)

	keyboard := keyboards.BuildCollectionsSortingKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleSortingCollectionsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackSortingCollectionsSelectBack:
		HandleCollectionsCommand(app, session)
		return

	case states.CallbackSortingCollectionsSelectAllReset:
		handleSortingCollectionsAllReset(app, session)
		return

	case states.CallbackSortingCollectionsSelectID:
		session.CollectionsState.Sorting.Field = "id"

	case states.CallbackSortingCollectionsSelectName:
		session.CollectionsState.Sorting.Field = "name"

	case states.CallbackSortingCollectionsSelectCreatedAt:
		session.CollectionsState.Sorting.Field = "created_at"

	case states.CallbackSortingCollectionsSelectTotalFilms:
		session.CollectionsState.Sorting.Field = "total_films"
	}

	handleSortingCollectionsDirection(app, session)
}

func HandleSortingCollectionsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleSortingCollectionsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessSortingCollectionsAwaitingDirection:
		parseSortingCollectionsDirection(app, session)
	}

}

func handleSortingCollectionsAllReset(app models.App, session *models.Session) {
	session.CollectionsState.Sorting.ResetSorting()

	msg := translator.Translate(session.Lang, "sortingResetSuccess", nil, nil)

	app.SendMessage(msg, nil)

	HandleSortingCollectionsCommand(app, session)
}

func handleSortingCollectionsDirection(app models.App, session *models.Session) {
	msg := messages.BuildSelectedSortMessage(session, session.CollectionsState.Sorting)

	keyboard := keyboards.BuildSortingDirectionKeyboard(session, session.CollectionsState.Sorting)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSortingCollectionsAwaitingDirection)
}

func parseSortingCollectionsDirection(app models.App, session *models.Session) {
	sorting := session.CollectionsState.Sorting

	if utils.IsReset(app.Upd) {
		sorting.Sort = ""
		handleSortingCollectionsReset(app, session)
		return
	}

	if utils.ParseCallback(app.Upd) == states.CallbacktDecrease {
		sorting.Direction = "-"
	}
	sorting.Sort = sorting.Direction + sorting.Field

	msg := translator.Translate(session.Lang, "sortingApplied", nil, nil)
	app.SendMessage(msg, nil)

	session.ClearAllStates()

	HandleSortingCollectionsCommand(app, session)
}

func handleSortingCollectionsReset(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "sortingResetSuccess", nil, nil)
	app.SendMessage(msg, nil)

	session.ClearAllStates()
	HandleSortingCollectionsCommand(app, session)
}
