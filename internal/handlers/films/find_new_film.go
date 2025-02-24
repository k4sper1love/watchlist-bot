package films

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
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
		app.SendMessage(messages.BuildFilmsNotFoundMessage(session), keyboards.BuildFindNewFilmsNotFoundKeyboard(session))
		clearStatesAndResetFilmsPage(session)
		return
	}

	app.SendMessage(messages.BuildFindNewFilmMessage(session, metadata), keyboards.BuildFindNewFilmKeyboard(session, metadata.CurrentPage, metadata.LastPage))
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

	case states.CallbackFindNewFilmNextPage, states.CallbackFindNewFilmPrevPage,
		states.CallbackFindNewFilmLastPage, states.CallbackFindNewFilmFirstPage:
		handleFindNewFilmPagination(app, session, callback)

	default:
		if strings.HasPrefix(callback, states.PrefixSelectFindNewFilm) {
			handleFindNewFilmSelect(app, session)
		}
	}
}

func handleFindNewFilmPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackFindNewFilmNextPage:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallbackFindNewFilmPrevPage:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallbackFindNewFilmLastPage:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallbackFindNewFilmFirstPage:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
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
		app.SendMessage(messages.BuildFilmsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFindNewFilmBack))
	} else {
		session.FilmDetailState.SetFromFilm(&session.FilmsState.Films[index])
		parseFilmImage(app, session, requestNewFilmComment)
	}
}

func getFilmsFromKinopoisk(session *models.Session) (*filters.Metadata, error) {
	films, metadata, err := parsing.GetFilmsFromKinopoisk(session)
	if err != nil {
		return nil, err
	}

	UpdateSessionWithFilms(session, films, metadata)
	return metadata, nil
}

func clearStatesAndResetFilmsPage(session *models.Session) {
	session.ClearAllStates()
	session.FilmsState.CurrentPage = 1
}
