package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleViewedFilmCommand(app models.App, session *models.Session) {
	session.FilmDetailState.IsViewed = true

	msg := "Фильм просмотрен✅\n"
	msg += "Поставьте оценку фильму\n\n"
	msg += "<i>Вы можете отменить \"просмотр\" фильма, нажав на кнопку отмены</i>"

	keyboard := builders.NewKeyboard(1).
		AddSkip().
		AddCancel().
		Build()

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessViewedFilmAwaitingRating)
}

func HandleViewedFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		session.FilmDetailState.Clear()
		HandleFilmsDetailCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessViewedFilmAwaitingRating:
		parseViewedFilmRating(app, session)

	case states.ProcessViewedFilmAwaitingReview:
		parseViewedFilmReview(app, session)
	}
}

func parseViewedFilmRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.UserRating = 0
	} else {
		session.FilmDetailState.UserRating = utils.ParseMessageFloat(app.Upd)
	}

	msg := "Оставьте отзыв к фильму"

	keyboard := builders.NewKeyboard(1).
		AddSkip().
		AddCancel().
		Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessViewedFilmAwaitingReview)
}

func parseViewedFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Review = ""
	} else {
		session.FilmDetailState.Review = utils.ParseMessageString(app.Upd)
	}

	finishUpdateFilmProcess(app, session)

	app.SendMessage("Успешно изменено!", nil)
	HandleFilmsDetailCommand(app, session)
}
