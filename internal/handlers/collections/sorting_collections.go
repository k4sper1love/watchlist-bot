package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

func HandleSortingCollectionsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceSorting(session), keyboards.BuildCollectionsSortingKeyboard(session))
}

func HandleSortingCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackSortingCollectionsBack:
		handleWithResetCollectionsPage(app, session, HandleCollectionsCommand)
	case states.CallbackSortingCollectionsAllReset:
		handleSortingCollectionsReset(app, session, HandleCollectionsCommand)
	default:
		handleSortingCollectionsSelect(app, session, callback)
	}
}

func handleSortingCollectionsSelect(app models.App, session *models.Session, callback string) {
	if strings.HasPrefix(callback, states.PrefixSortingCollectionsSelect) {
		session.CollectionsState.Sorting.Field = strings.TrimPrefix(callback, states.PrefixSortingCollectionsSelect)
		handleSortingCollectionsDirection(app, session)
	}
}

func HandleSortingCollectionsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		handleWithResetCollectionsPage(app, session, HandleSortingCollectionsCommand)
		return
	}

	switch session.State {
	case states.ProcessSortingCollectionsAwaitingDirection:
		parseSortingCollectionsDirection(app, session)
	}
}

func handleSortingCollectionsDirection(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestSortDirection(session, session.CollectionsState.Sorting), keyboards.BuildSortingDirectionKeyboard(session, session.CollectionsState.Sorting))
	session.SetState(states.ProcessSortingCollectionsAwaitingDirection)
}

func parseSortingCollectionsDirection(app models.App, session *models.Session) {
	if utils.IsReset(app.Update) {
		handleSortingCollectionsReset(app, session, HandleSortingCollectionsCommand)
		return
	}

	if utils.IsDecrease(app.Update) {
		session.CollectionsState.Sorting.Direction = "-"
	}

	session.CollectionsState.Sorting.SetSort()
	app.SendMessage(messages.SortingApplied(session, session.CollectionsState.Sorting), nil)
	handleWithResetCollectionsPage(app, session, HandleCollectionsCommand)
}

func handleSortingCollectionsReset(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.CollectionsState.Sorting.ResetSorting()
	app.SendMessage(messages.ResetSortingSuccess(session), nil)
	handleWithResetCollectionsPage(app, session, next)
}

func handleWithResetCollectionsPage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.ClearAllStates()
	session.CollectionsState.CurrentPage = 1
	next(app, session)
}
