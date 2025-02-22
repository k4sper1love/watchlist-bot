package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleViewedFilmCommand(app models.App, session *models.Session) {
	session.FilmDetailState.IsViewed = true

	part1 := translator.Translate(session.Lang, "viewedFilmRequestRating", nil, nil)
	part2 := translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)

	msg := fmt.Sprintf("✔️ <b>%s</b>\n\n<i>%s</i>", part1, part2)

	keyboard := keyboards.BuildFilmViewedKeyboard(session)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessViewedFilmAwaitingRating)
}

func HandleViewedFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
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
	if utils.IsSkip(app.Update) {
		session.FilmDetailState.UserRating = 0
	} else {
		session.FilmDetailState.UserRating = utils.ParseMessageFloat(app.Update)
	}

	part1 := translator.Translate(session.Lang, "viewedFilmRequestReview", nil, nil)
	part2 := translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)

	msg := fmt.Sprintf("✔️ <b>%s</b>\n\n<i>%s</i>", part1, part2)

	keyboard := keyboards.BuildFilmViewedKeyboard(session)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessViewedFilmAwaitingReview)
}

func parseViewedFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		session.FilmDetailState.Review = ""
	} else {
		session.FilmDetailState.Review = utils.ParseMessageString(app.Update)
	}

	finishUpdateFilmProcess(app, session, HandleFilmsDetailCommand)
}
