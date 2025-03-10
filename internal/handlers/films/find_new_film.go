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

func HandleFindNewFilmCommand(app models.App, session *models.Session) {
	metadata, err := getFilmsFromKinopoisk(session)
	if err != nil {
		handleKinopoiskError(app, session, err)
		clearStatesAndResetFilmsPage(session)
		return
	}

	if metadata.TotalRecords == 0 {
		app.SendMessage(messages.FilmsNotFound(session), keyboards.FindNewFilmsNotFound(session))
		clearStatesAndResetFilmsPage(session)
		return
	}

	app.SendMessage(messages.FindNewFilm(session, metadata), keyboards.FindNewFilm(session, metadata.CurrentPage, metadata.LastPage))
}

func HandleFindNewFilmButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackFindNewFilmBack:
		session.ClearAllStates()
		HandleFilmsCommand(app, session)

	case states.CallbackFindNewFilmAgain:
		session.ClearAllStates()
		handleNewFilmFind(app, session)

	case states.CallbackFindNewFilmPageNext, states.CallbackFindNewFilmPagePrev,
		states.CallbackFindNewFilmPageLast, states.CallbackFindNewFilmPageFirst:
		handleFindNewFilmPagination(app, session, callback)

	default:
		if strings.HasPrefix(callback, states.PrefixSelectFindNewFilm) {
			handleFindNewFilmSelect(app, session)
		}
	}
}

func handleFindNewFilmPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackFindNewFilmPageNext:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallbackFindNewFilmPagePrev:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallbackFindNewFilmPageLast:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallbackFindNewFilmPageFirst:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1
	}

	HandleFindNewFilmCommand(app, session)
}

func handleFindNewFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)
	indexStr := strings.TrimPrefix(callback, states.PrefixSelectFindNewFilm)

	if index, err := strconv.Atoi(indexStr); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallbackFindNewFilmBack))
	} else {
		session.FilmDetailState.SetFromFilm(&session.FilmsState.Films[index])
		parser.ParseFilmImageFromURL(app, session, session.FilmDetailState.ImageURL, requestNewFilmComment)
	}
}

func getFilmsFromKinopoisk(session *models.Session) (*filters.Metadata, error) {
	films, metadata, err := parsing.GetFilmsFromKinopoisk(session)
	if err != nil {
		return nil, err
	}

	updateSessionWithFilms(session, films, metadata)
	return metadata, nil
}

func clearStatesAndResetFilmsPage(session *models.Session) {
	session.ClearAllStates()
	session.FilmsState.CurrentPage = 1
}
