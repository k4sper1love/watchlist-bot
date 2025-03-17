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

// HandleAddCollectionToFilmCommand handles the command for adding a collection to a film.
// Retrieves collections excluding the current film and sends a paginated list with navigation buttons.
func HandleAddCollectionToFilmCommand(app models.App, session *models.Session) {
	if metadata, err := getCollectionsExcludeFilm(app, session); err != nil {
		app.SendMessage(messages.CollectionsFailure(session), nil)
	} else if metadata.TotalRecords == 0 {
		app.SendMessage(messages.CollectionsNotFound(session), keyboards.CollectionToFilmNotFound(session))
	} else {
		app.SendMessage(messages.ChoiceCollection(session), keyboards.AddCollectionToFilm(session))
	}
}

// HandleAddCollectionToFilmButtons handles button interactions related to adding a collection to a film.
// Supports actions like going back, searching, pagination, and selecting specific collections.
func HandleAddCollectionToFilmButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallAddCollectionToFilmBack:
		films.HandleFilmDetailCommand(app, session)

	case states.CallAddCollectionToFilmFind:
		handleAddCollectionToFilmFind(app, session)

	case states.CallAddCollectionToFilmAgain:
		session.CollectionsState.Name = ""
		handleAddCollectionToFilmFind(app, session)

	case states.CallAddCollectionToFilmReset:
		session.CollectionsState.Name = ""
		HandleAddCollectionToFilmCommand(app, session)

	default:
		if strings.HasPrefix(callback, states.AddCollectionToFilmPage) {
			handleAddCollectionToFilmPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectCFCollection) {
			handleAddCollectionToFilmSelect(app, session, callback)
		}
	}
}

// HandleAddCollectionToFilmProcess processes the workflow for adding a collection to a film.
// Handles states like awaiting a collection name input.
func HandleAddCollectionToFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleAddCollectionToFilmCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitAddCollectionToFilmName:
		parseAddCollectionToFilmName(app, session)
	}
}

// handleAddCollectionToFilmPagination processes pagination actions for the collection list.
// Updates the current page in the session and reloads the collection list.
func handleAddCollectionToFilmPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallAddCollectionToFilmPageNext:
		if session.CollectionFilmsState.CurrentPage >= session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage++

	case states.CallAddCollectionToFilmPagePrev:
		if session.CollectionFilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage--

	case states.CallAddCollectionToFilmPageLast:
		if session.CollectionFilmsState.CurrentPage == session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = session.CollectionFilmsState.LastPage

	case states.CallAddCollectionToFilmPageFirst:
		if session.CollectionFilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = 1
	}

	HandleAddCollectionToFilmCommand(app, session)
}

// handleAddCollectionToFilmSelect processes the selection of a collection from the list.
// Parses the collection ID and adds the film to the selected collection.
func handleAddCollectionToFilmSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectCFCollection)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.CollectionsFailure(session), keyboards.Back(session, states.CallCollectionFilmsFromFilm))
	} else {
		session.CollectionDetailState.Collection.ID = id
		addFilmToCollection(app, session)
	}
}

// handleAddCollectionToFilmFind prompts the user to enter the name of a collection to search for.
func handleAddCollectionToFilmFind(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionName(session), keyboards.Cancel(session))
	session.SetState(states.AwaitAddCollectionToFilmName)
}

// parseAddCollectionToFilmName processes the collection name provided by the user.
func parseAddCollectionToFilmName(app models.App, session *models.Session) {
	session.CollectionsState.Name = utils.ParseMessageString(app.Update)
	session.CollectionFilmsState.CurrentPage = 1

	session.ClearState()
	HandleAddCollectionToFilmCommand(app, session)
}

// getCollectionsExcludeFilm retrieves a paginated list of collections excluding the current film.
func getCollectionsExcludeFilm(app models.App, session *models.Session) (*filters.Metadata, error) {
	collectionsResponse, err := watchlist.GetCollectionsExcludeFilm(app, session)
	if err != nil {
		return nil, err
	}

	updateSessionWithCollectionsExcludeFilm(session, collectionsResponse)
	return &collectionsResponse.Metadata, nil
}

// updateSessionWithCollectionsExcludeFilm updates the session with the retrieved collections and their metadata.
func updateSessionWithCollectionsExcludeFilm(session *models.Session, collectionsResponse *models.CollectionsResponse) {
	session.CollectionsState.Collections = collectionsResponse.Collections
	session.CollectionFilmsState.CurrentPage = collectionsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = collectionsResponse.Metadata.LastPage
	session.CollectionFilmsState.TotalRecords = collectionsResponse.Metadata.TotalRecords
}
