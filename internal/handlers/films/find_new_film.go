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
	if metadata, err := getFilmsFromKinopoisk(session); err != nil {
		handleKinopoiskError(app, session, err)
		clearStatesAndResetFilmsPage(session)
	} else {
		app.SendMessage(messages.FindNewFilm(session, metadata), keyboards.FindNewFilm(session, metadata.CurrentPage, metadata.LastPage))
	}
}

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

func handleFindNewFilmSelect(app models.App, session *models.Session, callback string) {
	if index, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectNewFilm)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallFindNewFilmBack))
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
