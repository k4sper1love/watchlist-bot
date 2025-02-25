package collectionFilms

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

func HandleAddCollectionToFilmCommand(app models.App, session *models.Session) {
	if metadata, err := GetCollectionsExcludeFilm(app, session); err != nil {
		app.SendMessage(messages.BuildCollectionsFailureMessage(session), nil)
	} else if metadata.TotalRecords == 0 {
		app.SendMessage(messages.BuildCollectionsNotFoundMessage(session), keyboards.BuildAddCollectionToFilmNotFoundKeyboard(session))
	} else {
		app.SendMessage(messages.BuildChoiceCollectionMessage(session), keyboards.BuildAddCollectionToFilmKeyboard(session))
	}
}

func HandleAddCollectionToFilmButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackAddCollectionToFilmBack:
		films.HandleFilmsDetailCommand(app, session)

	case states.CallbackAddCollectionToFilmFind:
		handleAddCollectionToFilmFind(app, session)

	case states.CallbackAddCollectionToFilmAgain:
		session.CollectionsState.Name = ""
		handleAddCollectionToFilmFind(app, session)

	case states.CallbackAddCollectionToFilmReset:
		session.CollectionsState.Name = ""
		HandleAddCollectionToFilmCommand(app, session)

	case states.CallbackAddCollectionToFilmNextPage, states.CallbackAddCollectionToFilmPrevPage,
		states.CallbackAddCollectionToFilmLastPage, states.CallbackAddCollectionToFilmFirstPage:
		handleAddCollectionToFilmPagination(app, session, callback)

	default:
		if strings.HasPrefix(callback, states.PrefixSelectCFCollection) {
			handleAddCollectionToFilmSelect(app, session, callback)
		}
	}
}

func HandleAddCollectionToFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleAddCollectionToFilmCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAddCollectionToFilmAwaitingName:
		parseAddCollectionToFilmName(app, session)
	}
}

func handleAddCollectionToFilmPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackAddCollectionToFilmNextPage:
		if session.CollectionFilmsState.CurrentPage >= session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage++

	case states.CallbackAddCollectionToFilmPrevPage:
		if session.CollectionFilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage--

	case states.CallbackAddCollectionToFilmLastPage:
		if session.CollectionFilmsState.CurrentPage == session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = session.CollectionFilmsState.LastPage

	case states.CallbackAddCollectionToFilmFirstPage:
		if session.CollectionFilmsState.CurrentPage == 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = 1
	}

	HandleAddCollectionToFilmCommand(app, session)
}

func handleAddCollectionToFilmSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, states.PrefixSelectCFCollection)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.BuildCollectionsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackCollectionFilmsFromFilm))
	} else {
		session.CollectionDetailState.Collection.ID = id
		addFilmToCollection(app, session)
	}
}

func handleAddCollectionToFilmFind(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildCollectionRequestNameMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessAddCollectionToFilmAwaitingName)
}

func parseAddCollectionToFilmName(app models.App, session *models.Session) {
	session.CollectionsState.Name = utils.ParseMessageString(app.Update)
	session.CollectionFilmsState.CurrentPage = 1

	session.ClearState()
	HandleAddCollectionToFilmCommand(app, session)
}

func GetCollectionsExcludeFilm(app models.App, session *models.Session) (*filters.Metadata, error) {
	collectionsResponse, err := watchlist.GetCollectionsExcludeFilm(app, session)
	if err != nil {
		return nil, err
	}

	updateSessionWithCollectionsExcludeFilm(session, collectionsResponse)
	return &collectionsResponse.Metadata, nil
}

func updateSessionWithCollectionsExcludeFilm(session *models.Session, collectionsResponse *models.CollectionsResponse) {
	session.CollectionsState.Collections = collectionsResponse.Collections
	session.CollectionFilmsState.CurrentPage = collectionsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = collectionsResponse.Metadata.LastPage
	session.CollectionFilmsState.TotalRecords = collectionsResponse.Metadata.TotalRecords
}
