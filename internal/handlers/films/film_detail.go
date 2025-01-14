package films

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleFilmsDetailCommand(app models.App, session *models.Session) {
	index := session.FilmDetailState.Index
	films := session.FilmsState.Films
	film := films[index]

	session.FilmDetailState.Film = film

	msg := messages.BuildFilmDetailMessage(session, &film)

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
		session.FilmsState.CurrentPage = 1
		HandleFilmsCommand(app, session)

	case states.CallbackFilmDetailViewed:
		HandleViewedFilmCommand(app, session)

	case states.CallbackFilmDetailFavorite:
		handleFavoriteFilm(app, session)
	}
}

func handleFavoriteFilm(app models.App, session *models.Session) {
	session.FilmDetailState.IsFavorite = !session.FilmDetailState.Film.IsFavorite

	finishUpdateFilmProcess(app, session, HandleFilmsDetailCommand)
}

func getFilmsLastIndex(session *models.Session) int {
	return len(session.FilmsState.Films) - 1
}

func UpdateFilmInList(app models.App, session *models.Session) error {
	film, err := watchlist.GetFilm(app, session)
	if err != nil {
		return err
	}

	UpdateSessionWithFilm(session, film)

	return nil
}

func UpdateSessionWithFilm(session *models.Session, film *apiModels.Film) {
	index := session.FilmDetailState.Index
	session.FilmsState.Films[index] = *film
}
