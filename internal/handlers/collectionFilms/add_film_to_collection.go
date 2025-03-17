package collectionFilms

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

// HandleAddFilmToCollectionCommand handles the command for adding a film to a collection.
// Retrieves films excluding the current collection and sends a paginated list with navigation buttons.
func HandleAddFilmToCollectionCommand(app models.App, session *models.Session) {
	if metadata, err := getFilmsExcludeCollection(app, session); err != nil {
		app.SendMessage(messages.FilmsFailure(session), nil)
	} else if metadata.TotalRecords == 0 {
		app.SendMessage(messages.FilmsNotFound(session), keyboards.FilmToCollectionNotFound(session))
	} else {
		app.SendMessage(messages.ChoiceFilm(session), keyboards.AddFilmToCollection(session))
	}
}

// HandleAddFilmToCollectionButtons handles button interactions related to adding a film to a collection.
// Supports actions like going back, searching, pagination, and selecting specific films.
func HandleAddFilmToCollectionButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch utils.ParseCallback(app.Update) {
	case states.CallAddFilmToCollectionBack:
		HandleOptionsFilmToCollectionCommand(app, session)

	case states.CallAddFilmToCollectionFind:
		handleAddFilmToCollectionFind(app, session)

	case states.CallAddFilmToCollectionAgain:
		session.FilmsState.Title = ""
		handleAddFilmToCollectionFind(app, session)

	case states.CallAddFilmToCollectionReset:
		session.FilmsState.Title = ""
		HandleAddFilmToCollectionCommand(app, session)

	default:
		if strings.HasPrefix(callback, states.AddFilmToCollectionPage) {
			handleAddFilmToCollectionPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectCFFilm) {
			handleAddFilmToCollectionSelect(app, session, callback)
		}
	}
}

// HandleAddFilmToCollectionProcess processes the workflow for adding a film to a collection.
// Handles states like awaiting a film title input.
func HandleAddFilmToCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleAddFilmToCollectionCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitAddFilmToCollectionTitle:
		parseAddFilmToCollectionTitle(app, session)
	}
}

// handleAddFilmToCollectionPagination processes pagination actions for the film list.
// Updates the current page in the session and reloads the film list.
func handleAddFilmToCollectionPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallAddFilmToCollectionPageNext:
		if session.CollectionFilmsState.CurrentPage >= session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage++

	case states.CallAddFilmToCollectionPagePrev:
		if session.CollectionFilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage--

	case states.CallAddFilmToCollectionPageLast:
		if session.CollectionFilmsState.CurrentPage == session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = session.CollectionFilmsState.LastPage

	case states.CallAddFilmToCollectionPageFirst:
		if session.CollectionFilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = 1
	}

	HandleAddFilmToCollectionCommand(app, session)
}

// handleAddFilmToCollectionSelect processes the selection of a film from the list.
// Parses the film ID and adds the selected film to the collection.
func handleAddFilmToCollectionSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectCFFilm)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallFilmToCollectionOptionExisting))
	} else {
		session.FilmDetailState.Film.ID = id
		AddFilmToCollection(app, session)
	}
}

// handleAddFilmToCollectionFind prompts the user to enter the title of a film to search for.
func handleAddFilmToCollectionFind(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.AwaitAddFilmToCollectionTitle)
}

// parseAddFilmToCollectionTitle processes the film title provided by the user.
func parseAddFilmToCollectionTitle(app models.App, session *models.Session) {
	session.FilmsState.Title = utils.ParseMessageString(app.Update)
	session.CollectionFilmsState.CurrentPage = 1

	session.ClearState()
	HandleAddFilmToCollectionCommand(app, session)
}

// getFilmsExcludeCollection retrieves a paginated list of films excluding the current collection.
func getFilmsExcludeCollection(app models.App, session *models.Session) (*filters.Metadata, error) {
	filmsResponse, err := watchlist.GetFilmsExcludeCollection(app, session)
	if err != nil {
		return nil, err
	}

	updateSessionWithFilmsExcludeCollection(session, filmsResponse)
	return &filmsResponse.Metadata, nil
}

// updateSessionWithFilmsExcludeCollection updates the session with the retrieved films and their metadata.
func updateSessionWithFilmsExcludeCollection(session *models.Session, filmsResponse *models.FilmsResponse) {
	session.FilmsState.Films = filmsResponse.Films
	session.CollectionFilmsState.CurrentPage = filmsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = filmsResponse.Metadata.LastPage
}
