package collectionFilms

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
	"strconv"
	"strings"
)

func HandleAddFilmToCollectionCommand(app models.App, session *models.Session) {
	films, err := GetFilmsExcludeCollection(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	if len(films) == 0 {
		msg := "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
		keyboard := keyboards.BuildAddFilmToCollectionNotFoundKeyboard(session)
		app.SendMessage(msg, keyboard)
		return
	}

	choiceMsg := translator.Translate(session.Lang, "choiceFilm", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>", choiceMsg)
	keyboard := keyboards.BuildAddFilmToCollectionKeyboard(session)
	app.SendMessage(msg, keyboard)
}

func HandleAddFilmToCollectionButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	switch {
	case callback == states.CallbackAddFilmToCollectionBack:
		HandleOptionsFilmToCollectionCommand(app, session)

	case strings.HasPrefix(callback, "select_cf_film_"):
		HandleAddFilmToCollectionSelect(app, session)

	case callback == states.CallbackAddFilmToCollectionFind:
		handleAddFilmToCollectionFind(app, session)

	case callback == states.CallbackAddFilmToCollectionAgain:
		session.FilmsState.Title = ""
		handleAddFilmToCollectionFind(app, session)

	case callback == states.CallbackAddFilmToCollectionReset:
		session.FilmsState.Title = ""
		HandleAddFilmToCollectionCommand(app, session)

	case callback == states.CallbackAddFilmToCollectionNextPage:
		if session.CollectionFilmsState.CurrentPage < session.CollectionFilmsState.LastPage {
			session.CollectionFilmsState.CurrentPage++
			HandleAddFilmToCollectionCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAddFilmToCollectionPrevPage:
		if session.CollectionFilmsState.CurrentPage > 1 {
			session.CollectionFilmsState.CurrentPage--
			HandleAddFilmToCollectionCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAddFilmToCollectionLastPage:
		if session.CollectionFilmsState.CurrentPage != session.CollectionFilmsState.LastPage {
			session.CollectionFilmsState.CurrentPage = session.CollectionFilmsState.LastPage
			HandleAddFilmToCollectionCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAddFilmToCollectionFirstPage:
		if session.CollectionFilmsState.CurrentPage != 1 {
			session.CollectionFilmsState.CurrentPage = 1
			HandleAddFilmToCollectionCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}
	}
}

func HandleAddFilmToCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleAddFilmToCollectionCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAddFilmToCollectionAwaitingTitle:
		parseAddFilmToCollectionTitle(app, session)
	}
}

func HandleAddFilmToCollectionSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_cf_film_")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "getFilmFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackOptionsFilmToCollectionExisting).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		sl.Log.Error("failed to parse film ID", slog.Any("error", err), slog.String("callback", callback))
		return
	}

	session.FilmDetailState.Film.ID = id

	addFilmToCollection(app, session)
}

func handleAddFilmToCollectionFind(app models.App, session *models.Session) {
	msg := "❓ " + translator.Translate(session.Lang, "filmRequestTitle", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAddFilmToCollectionAwaitingTitle)
}

func parseAddFilmToCollectionTitle(app models.App, session *models.Session) {
	title := utils.ParseMessageString(app.Upd)

	session.FilmsState.Title = title
	session.CollectionFilmsState.CurrentPage = 1

	session.ClearState()

	HandleAddFilmToCollectionCommand(app, session)
}

func GetFilmsExcludeCollection(app models.App, session *models.Session) ([]apiModels.Film, error) {
	filmsResponse, err := watchlist.GetFilmsExcludeCollection(app, session)
	if err != nil {
		return nil, err
	}

	session.FilmsState.Films = filmsResponse.Films
	session.CollectionFilmsState.CurrentPage = filmsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = filmsResponse.Metadata.LastPage

	return filmsResponse.Films, nil
}
