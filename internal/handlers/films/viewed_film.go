package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleViewedFilmCommand(app models.App, session *models.Session) {
	session.FilmDetailState.SetViewed(true)
	requestViewedFilmUserRating(app, session)
}

func HandleViewedFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleFilmsDetailCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessViewedFilmAwaitingUserRating:
		parseViewedFilmUserRating(app, session)
	case states.ProcessViewedFilmAwaitingReview:
		parseViewedFilmReview(app, session)
	}
}

func requestViewedFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildViewedFilmRequestUserRatingMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessViewedFilmAwaitingUserRating)
}

func parseViewedFilmUserRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestViewedFilmReview(app, session)
		return
	}

	if userRating, ok := parseAndValidateNumber(app, session, 1, 10, utils.ParseMessageFloat); !ok {
		requestViewedFilmUserRating(app, session)
	} else {
		session.FilmDetailState.UserRating = userRating
		requestViewedFilmReview(app, session)
	}
}

func requestViewedFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildViewedFilmRequestReviewMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessViewedFilmAwaitingReview)
}

func parseViewedFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		finishUpdateFilmProcess(app, session, HandleFilmsDetailCommand)
		return
	}

	if review, ok := parseAndValidateString(app, session, 0, 500); !ok {
		requestViewedFilmReview(app, session)
	} else {
		session.FilmDetailState.Review = review
		finishUpdateFilmProcess(app, session, HandleFilmsDetailCommand)
	}
}
