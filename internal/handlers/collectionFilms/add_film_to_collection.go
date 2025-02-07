package collectionFilms

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log"
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
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAddFilmToCollectionBack).Build(session.Lang)
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

func HandleAddFilmToCollectionSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_cf_film_")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "getFilmFailure", nil, nil)
		app.SendMessage(msg, nil)
		log.Printf("error parsing film ID: %v", err)
		return
	}

	session.FilmDetailState.Film.ID = id

	addFilmToCollection(app, session)
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
