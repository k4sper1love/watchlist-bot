package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"strconv"
	"strings"
)

func HandleFilmsCommand(app models.App, session *models.Session) {
	filmsResponse, err := getFilms(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	metadata := filmsResponse.Metadata

	msg := builders.BuildFilmsMessage(filmsResponse)

	keyboard := builders.NewKeyboard(1).
		AddFilmsSelect(filmsResponse).
		AddFilmsNew().
		AddNavigation(metadata.CurrentPage, metadata.LastPage, states.CallbackFilmsPrevPage, states.CallbackFilmsNextPage).
		AddBack(states.CallbackFilmsBack).
		Build()

	app.SendMessage(msg, keyboard)
}

func HandleFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackFilmsBack:
		general.HandleMenuCommand(app, session)

	case callback == states.CallbackFilmsNextPage:
		if session.FilmsState.CurrentPage < session.FilmsState.LastPage {
			session.FilmsState.CurrentPage++
			HandleFilmsCommand(app, session)
		} else {
			app.SendMessage("Вы уже на последней странице", nil)
		}

	case callback == states.CallbackFilmsPrevPage:
		if session.FilmsState.CurrentPage > 1 {
			session.FilmsState.CurrentPage--
			HandleFilmsCommand(app, session)
		} else {
			app.SendMessage("Вы уже на первой странице", nil)
		}

	case callback == states.CallbackFilmsNew:
		HandleNewFilmCommand(app, session)

	case callback == states.CallbackFilmsManage:
		HandleManageFilmCommand(app, session)

	case strings.HasPrefix(callback, "select_film_"):
		handleFilmSelect(app, session)
	}
}

func handleFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	indexStr := strings.TrimPrefix(callback, "select_film_")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		app.SendMessage("Ошибка при получении ID фильма.", nil)
		log.Printf("error parsing film index: %v", err)
		return
	}

	session.FilmDetailState.Index = index
	HandleFilmsDetailCommand(app, session)
}

func getFilms(app models.App, session *models.Session) (*models.FilmsResponse, error) {
	filmsResponse, err := watchlist.GetFilms(app, session)
	if err != nil {
		return nil, err
	}

	session.FilmsState.Object = filmsResponse.Films
	session.FilmsState.LastPage = filmsResponse.Metadata.LastPage
	session.FilmsState.TotalRecords = filmsResponse.Metadata.TotalRecords

	return filmsResponse, nil
}

func updateFilmsList(app models.App, session *models.Session, next bool) error {
	currentPage := session.FilmsState.CurrentPage
	lastPage := session.FilmsState.LastPage

	switch next {
	case true:
		if currentPage < lastPage {
			session.FilmsState.CurrentPage++
		} else {
			return fmt.Errorf("Вы уже на последней странице")
		}
	case false:
		if currentPage > 1 {
			session.FilmsState.CurrentPage--
		} else {
			return fmt.Errorf("Вы уже на первой странице")
		}
	}

	_, err := getFilms(app, session)
	if err != nil {
		return err
	}

	return nil
}
