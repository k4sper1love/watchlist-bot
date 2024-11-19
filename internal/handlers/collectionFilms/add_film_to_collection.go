package collectionFilms

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"strconv"
	"strings"
)

func HandleAddFilmToCollectionCommand(app models.App, session *models.Session) {
	if err := GetFilmsExcludeCollection(app, session); err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	msg := "Выберите, какой фильм добавить в коллекцию?"

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
			app.SendMessage("Вы уже на последней странице", nil)
		}

	case callback == states.CallbackAddFilmToCollectionPrevPage:
		if session.CollectionFilmsState.CurrentPage > 1 {
			session.CollectionFilmsState.CurrentPage--
			HandleAddFilmToCollectionCommand(app, session)
		} else {
			app.SendMessage("Вы уже на первой странице", nil)
		}
	}
}

func HandleAddFilmToCollectionSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_cf_film_")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		app.SendMessage("Ошибка при получении ID фильма.", nil)
		log.Printf("error parsing film ID: %v", err)
		return
	}

	session.FilmDetailState.Film.ID = id

	addFilmToCollection(app, session)
}

func GetFilmsExcludeCollection(app models.App, session *models.Session) error {
	filmsResponse, err := watchlist.GetFilmsExcludeCollection(app, session)
	if err != nil {
		return err
	}

	session.FilmsState.Films = filmsResponse.Films
	session.CollectionFilmsState.CurrentPage = filmsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = filmsResponse.Metadata.LastPage

	return nil
}
