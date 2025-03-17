package films

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/parsing"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

// HandleFindNewFilmCommand handles the command for searching new films from Kinopoisk.
// Retrieves paginated films and sends a message with their details and navigation buttons.
func HandleFindNewFilmCommand(app models.App, session *models.Session) {
	if metadata, err := getFilmsFromKinopoisk(session); err != nil {
		handleKinopoiskError(app, session, err)
		clearStatesAndResetFilmsPage(session)
	} else {
		app.SendMessage(messages.FindNewFilm(session, metadata), keyboards.FindNewFilm(session, metadata.CurrentPage, metadata.LastPage))
	}
}

// HandleFindNewFilmButtons handles button interactions related to the search results of new films.
// Supports actions like going back, refreshing the search, pagination, and selecting specific films.
func HandleFindNewFilmButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFindNewFilmBack:
		session.ClearAllStates()
		HandleFilmsCommand(app, session)

	case states.CallFindNewFilmAgain:
		session.ClearAllStates()
		handleNewFilmFind(app, session)

	default:
		if strings.HasPrefix(callback, states.FindNewFilmPage) {
			handleFindNewFilmPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectNewFilm) {
			handleFindNewFilmSelect(app, session, callback)
		}
	}
}

// handleFindNewFilmPagination processes pagination actions for the search results of new films.
// Updates the current page in the session and reloads the films list.
func handleFindNewFilmPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallFindNewFilmPageNext:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallFindNewFilmPagePrev:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallFindNewFilmPageLast:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallFindNewFilmPageFirst:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1
	}

	HandleFindNewFilmCommand(app, session)
}

// handleFindNewFilmSelect processes the selection of a film from the search results.
// Parses the film index and navigates to the detailed view of the selected film.
func handleFindNewFilmSelect(app models.App, session *models.Session, callback string) {
	if index, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectNewFilm)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallFindNewFilmBack))
	} else {
		session.FilmDetailState.SetFromFilm(&session.FilmsState.Films[index])
		parser.ParseFilmImageFromURL(app, session, session.FilmDetailState.ImageURL, requestNewFilmComment)
	}
}

// getFilmsFromKinopoisk retrieves a paginated list of films from Kinopoisk using the Parsing service.
// Updates the session with the retrieved films and their metadata.
func getFilmsFromKinopoisk(session *models.Session) (*filters.Metadata, error) {
	films, metadata, err := parsing.GetFilmsFromKinopoisk(session)
	if err != nil {
		return nil, err
	}

	updateSessionWithFilms(session, films, metadata)
	return metadata, nil
}

// clearStatesAndResetFilmsPage clears all session states and resets the current page of the films list.
func clearStatesAndResetFilmsPage(session *models.Session) {
	session.ClearAllStates()
	session.FilmsState.CurrentPage = 1
}
