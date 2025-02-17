package films

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/parsing"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log"
	"strconv"
	"strings"
)

func HandleFindNewFilmCommand(app models.App, session *models.Session) {
	metadata, err := getFilmsFromKinopoisk(session)
	if err != nil {
		session.FilmsState.CurrentPage = 1
		handleKinopoiskError(app, session, err)
		return
	}

	if metadata.TotalRecords == 0 {
		msg := "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
		keyboard := keyboards.NewKeyboard().AddAgain(states.CallbackFindNewFilmAgain).AddBack(states.CallbackFindNewFilmBack).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		session.ClearAllStates()
		session.FilmsState.CurrentPage = 1
		return
	}

	msg := messages.BuildFindNewFilmMessage(session, metadata)

	keyboard := keyboards.BuildFindNewFilmKeyboard(session, metadata.CurrentPage, metadata.LastPage)

	app.SendMessage(msg, keyboard)

}

func HandleFindNewFilmButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackFindNewFilmBack:
		session.ClearAllStates()
		HandleFilmsCommand(app, session)
		return

	case callback == states.CallbackFindNewFilmAgain:
		session.ClearAllStates()
		handleNewFilmFind(app, session)
		return

	case callback == states.CallbackFindNewFilmNextPage:
		if session.FilmsState.CurrentPage < session.FilmsState.LastPage {
			session.FilmsState.CurrentPage++
			HandleFindNewFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFindNewFilmPrevPage:
		if session.FilmsState.CurrentPage > 1 {
			session.FilmsState.CurrentPage--
			HandleFindNewFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFindNewFilmLastPage:
		if session.FilmsState.CurrentPage != session.FilmsState.LastPage {
			session.FilmsState.CurrentPage = session.FilmsState.LastPage
			HandleFindNewFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFindNewFilmFirstPage:
		if session.FilmsState.CurrentPage != 1 {
			session.FilmsState.CurrentPage = 1
			HandleFindNewFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}
	case strings.HasPrefix(callback, "select_find_new_film_"):
		handleFindNewFilmSelect(app, session)
	}
}

func handleFindNewFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	indexStr := strings.TrimPrefix(callback, "select_find_new_film_")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		msg := "🚨️ " + translator.Translate(session.Lang, "getFilmFailure", nil, nil)
		app.SendMessage(msg, nil)
		log.Printf("error parsing find_new_film index: %v", err)
		return
	}

	film := session.FilmsState.Films[index]

	session.FilmDetailState.SetFromFilm(&film)

	imageURL, err := parseAndUploadImageFromURL(app, film.ImageURL)
	if err != nil {
		msg := "⚠️ " + translator.Translate(session.Lang, "getImageFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.FilmDetailState.SetImageURL("")
		requestNewFilmComment(app, session)
		return
	}
	session.FilmDetailState.SetImageURL(imageURL)

	handleFindNewFilmAdd(app, session)
}

func handleFindNewFilmAdd(app models.App, session *models.Session) {
	session.FilmsState.CurrentPage = 1
	session.FilmsState.Title = ""

	_, err := GetFilms(app, session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		session.FilmsState.CurrentPage = 1
		HandleFilmsCommand(app, session)
		return
	}

	requestNewFilmComment(app, session)
}

func getFilmsFromKinopoisk(session *models.Session) (*filters.Metadata, error) {
	films, metadata, err := parsing.GetFilmsFromKinopoisk(session)
	if err != nil {
		return nil, err
	}

	UpdateSessionWithFilms(session, films, metadata)

	return metadata, nil
}
