package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleViewedCollectionFilmCommand(app models.App, session *models.Session) {
	session.CollectionFilmState.IsViewed = true

	msg := "Фильм просмотрен✅\n"
	msg += "Поставьте оценку фильму\n\n"
	msg += "<i>Вы можете отменить \"просмотр\" фильма, нажав на кнопку отмены</i>"

	keyboard := builders.NewKeyboard(1).
		AddSkip().
		AddCancel().
		Build()

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessViewedCollectionFilmAwaitingRating)
}

func HandleViewedCollectionFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		session.CollectionFilmState.Clear()
		HandleCollectionFilmsDetailCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessViewedCollectionFilmAwaitingRating:
		parseViewedCollectionFilmRating(app, session)

	case states.ProcessViewedCollectionFilmAwaitingReview:
		parseViewedCollectionFilmReview(app, session)
	}
}

func parseViewedCollectionFilmRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.UserRating = 0
	} else {
		session.CollectionFilmState.UserRating = utils.ParseMessageFloat(app.Upd)
	}

	msg := "Оставьте отзыв к фильму"

	keyboard := builders.NewKeyboard(1).
		AddSkip().
		AddCancel().
		Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessViewedCollectionFilmAwaitingReview)
}

func parseViewedCollectionFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.Review = ""
	} else {
		session.CollectionFilmState.Review = utils.ParseMessageString(app.Upd)
	}

	finishUpdateCollectionFilmProcess(app, session)

	app.SendMessage("Успешно изменено!", nil)
	HandleCollectionFilmsDetailCommand(app, session)
}
