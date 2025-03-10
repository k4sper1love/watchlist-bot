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

func HandleAddFilmToCollectionCommand(app models.App, session *models.Session) {
	if metadata, err := getFilmsExcludeCollection(app, session); err != nil {
		app.SendMessage(messages.FilmsFailure(session), nil)
	} else if metadata.TotalRecords == 0 {
		app.SendMessage(messages.FilmsNotFound(session), keyboards.FilmToCollectionNotFound(session))
	} else {
		app.SendMessage(messages.ChoiceFilm(session), keyboards.AddFilmToCollection(session))
	}
}

func HandleAddFilmToCollectionButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch utils.ParseCallback(app.Update) {
	case states.CallbackAddFilmToCollectionBack:
		HandleOptionsFilmToCollectionCommand(app, session)

	case states.CallbackAddFilmToCollectionFind:
		handleAddFilmToCollectionFind(app, session)

	case states.CallbackAddFilmToCollectionAgain:
		session.FilmsState.Title = ""
		handleAddFilmToCollectionFind(app, session)

	case states.CallbackAddFilmToCollectionReset:
		session.FilmsState.Title = ""
		HandleAddFilmToCollectionCommand(app, session)

	case states.CallbackAddFilmToCollectionPageNext, states.CallbackAddFilmToCollectionPagePrev,
		states.CallbackAddFilmToCollectionPageLast, states.CallbackAddFilmToCollectionPageFirst:
		handleAddFilmToCollectionPagination(app, session, callback)

	default:
		if strings.HasPrefix(callback, states.PrefixSelectCFFilm) {
			handleAddFilmToCollectionSelect(app, session, callback)
		}
	}
}

func HandleAddFilmToCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleAddFilmToCollectionCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAddFilmToCollectionAwaitingTitle:
		parseAddFilmToCollectionTitle(app, session)
	}
}

func handleAddFilmToCollectionPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackAddFilmToCollectionPageNext:
		if session.CollectionFilmsState.CurrentPage >= session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage++

	case states.CallbackAddFilmToCollectionPagePrev:
		if session.CollectionFilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage--

	case states.CallbackAddFilmToCollectionPageLast:
		if session.CollectionFilmsState.CurrentPage == session.CollectionFilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = session.CollectionFilmsState.LastPage

	case states.CallbackAddFilmToCollectionPageFirst:
		if session.CollectionFilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionFilmsState.CurrentPage = 1
	}

	HandleAddFilmToCollectionCommand(app, session)
}

func handleAddFilmToCollectionSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, states.PrefixSelectCFFilm)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallbackOptionsFilmToCollectionExisting))
	} else {
		session.FilmDetailState.Film.ID = id
		addFilmToCollection(app, session)
	}
}

func handleAddFilmToCollectionFind(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.ProcessAddFilmToCollectionAwaitingTitle)
}

func parseAddFilmToCollectionTitle(app models.App, session *models.Session) {
	session.FilmsState.Title = utils.ParseMessageString(app.Update)
	session.CollectionFilmsState.CurrentPage = 1

	session.ClearState()
	HandleAddFilmToCollectionCommand(app, session)
}

func getFilmsExcludeCollection(app models.App, session *models.Session) (*filters.Metadata, error) {
	filmsResponse, err := watchlist.GetFilmsExcludeCollection(app, session)
	if err != nil {
		return nil, err
	}

	updateSessionWithFilmsExcludeCollection(session, filmsResponse)
	return &filmsResponse.Metadata, nil
}

func updateSessionWithFilmsExcludeCollection(session *models.Session, filmsResponse *models.FilmsResponse) {
	session.FilmsState.Films = filmsResponse.Films
	session.CollectionFilmsState.CurrentPage = filmsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = filmsResponse.Metadata.LastPage
}
