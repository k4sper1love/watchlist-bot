package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleFilmsDetailCommand(app models.App, session *models.Session) {
	index := session.FilmDetailState.Index
	films := session.FilmsState.Films
	film := films[index]

	if len(films) == 0 {
		app.SendMessage("Фильмы не найдены. Начните с начала", nil)
		return
	}

	if index == -1 || index >= len(films) {
		app.SendMessage("Неизвестный фильм. Начните с начала", nil)
		return
	}

	session.FilmDetailState.Film = film

	itemID := utils.GetItemID(index, session.FilmsState.CurrentPage, session.FilmsState.PageSize)

	msg := messages.BuildFilmDetailWithNumberMessage(itemID, &film)

	keyboard := keyboards.BuildFilmDetailKeyboard(session)

	imageURL := film.ImageURL

	app.SendImage(imageURL, msg, keyboard)
}

func HandleFilmsDetailButtons(app models.App, session *models.Session) {
	currentIndex := session.FilmDetailState.Index
	lastIndex := getFilmsLastIndex(session)

	switch utils.ParseCallback(app.Upd) {
	case states.CallbackFilmDetailNextPage:
		if currentIndex < lastIndex {
			session.FilmDetailState.Index++
			HandleFilmsDetailCommand(app, session)
		} else {
			if err := UpdateFilmsList(app, session, true); err != nil {
				app.SendMessage(err.Error(), nil)
				return
			}
			session.FilmDetailState.Index = 0
			HandleFilmsDetailCommand(app, session)
		}

	case states.CallbackFilmDetailPrevPage:
		if currentIndex > 0 {
			session.FilmDetailState.Index--
			HandleFilmsDetailCommand(app, session)
		} else {
			if err := UpdateFilmsList(app, session, false); err != nil {
				app.SendMessage(err.Error(), nil)
				return
			}

			session.FilmDetailState.Index = getFilmsLastIndex(session)
			HandleFilmsDetailCommand(app, session)
		}

	case states.CallbackFilmDetailBack:
		HandleFilmsCommand(app, session)

	case states.CallbackFilmDetailViewed:
		HandleViewedFilmCommand(app, session)
	}
}

func getFilmsLastIndex(session *models.Session) int {
	return len(session.FilmsState.Films) - 1
}
