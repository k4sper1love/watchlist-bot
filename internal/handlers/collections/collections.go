package collections

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

// HandleCollectionsCommand handles the command for listing collections.
// Retrieves paginated collections and sends a message with their details and navigation buttons.
func HandleCollectionsCommand(app models.App, session *models.Session) {
	// Clears the name used for finding collections in other contexts to ensure a fresh state.
	session.CollectionsState.Clear()

	if metadata, err := getCollections(app, session); err != nil {
		app.SendMessage(messages.CollectionsFailure(session), keyboards.Back(session, ""))
	} else {
		app.SendMessage(messages.Collections(session, metadata, false), keyboards.Collections(session, metadata.CurrentPage, metadata.LastPage))
	}
}

// HandleCollectionsButtons handles button interactions related to collections.
// Supports actions like going back, creating new collections, managing existing ones, searching, sorting, and pagination.
func HandleCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallCollectionsBack:
		general.HandleMenuCommand(app, session)

	case states.CallCollectionsNew:
		HandleNewCollectionCommand(app, session)

	case states.CallCollectionsManage:
		HandleManageCollectionCommand(app, session)

	case states.CallCollectionsFind:
		handleCollectionsFindByName(app, session)

	case states.CallCollectionsSorting:
		HandleSortingCollectionsCommand(app, session)

	case states.CallCollectionsFavorite:
		handleFavoriteCollection(app, session)

	default:
		if strings.HasPrefix(callback, states.CollectionsPage) {
			handleCollectionsPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectCollection) {
			handleCollectionSelect(app, session, callback)
		}
	}
}

// HandleCollectionProcess processes workflows related to collections.
// Handles states like awaiting a collection name input for search.
func HandleCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitCollectionsName:
		parser.ParseCollectionFindName(app, session, HandleFindCollectionsCommand)
	}
}

// handleCollectionsPagination processes pagination actions for the collections list.
// Updates the current page in the session and reloads the collections list.
func handleCollectionsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallCollectionsPageNext:
		if session.CollectionsState.CurrentPage >= session.CollectionsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage++

	case states.CallCollectionsPagePrev:
		if session.CollectionsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage--

	case states.CallCollectionsPageLast:
		if session.CollectionsState.CurrentPage == session.CollectionsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = session.CollectionsState.LastPage

	case states.CallCollectionsPageFirst:
		if session.CollectionsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = 1
	}

	HandleCollectionsCommand(app, session)
}

// handleCollectionSelect processes the selection of a collection from the list.
// Parses the collection ID and navigates to the films associated with the selected collection.
func handleCollectionSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectCollection)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.CollectionsFailure(session), keyboards.Back(session, states.CallCollectionsBack))
	} else {
		session.CollectionDetailState.ObjectID = id
		setContextAndHandleFilms(app, session)
	}
}

// handleCollectionsFindByName prompts the user to enter the name of a collection to search for.
func handleCollectionsFindByName(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionName(session), keyboards.Cancel(session))
	session.SetState(states.AwaitCollectionsName)
}

// handleFavoriteCollection toggles the favorite status of the current collection and updates it.
func handleFavoriteCollection(app models.App, session *models.Session) {
	session.CollectionDetailState.SetFavorite(!session.CollectionDetailState.Collection.IsFavorite)
	HandleUpdateCollection(app, session, films.HandleFilmsCommand)
}

// getCollections retrieves a paginated list of collections using the Watchlist service.
func getCollections(app models.App, session *models.Session) (*filters.Metadata, error) {
	collectionsResponse, err := watchlist.GetCollections(app, session)
	if err != nil {
		return nil, err
	}

	updateSessionWithCollections(session, collectionsResponse.Collections, &collectionsResponse.Metadata)
	return &collectionsResponse.Metadata, nil
}

// updateSessionWithCollections updates the session with the retrieved collections and their metadata.
func updateSessionWithCollections(session *models.Session, collections []apiModels.Collection, metadata *filters.Metadata) {
	session.CollectionsState.Collections = collections
	session.CollectionsState.LastPage = metadata.LastPage
}

// setContextAndHandleFilms sets the context to "collection" and navigates to the films associated with the collection.
func setContextAndHandleFilms(app models.App, session *models.Session) {
	session.SetContext(states.CtxCollection)
	session.FilmsState.CurrentPage = 1
	films.HandleFilmsCommand(app, session)
}
