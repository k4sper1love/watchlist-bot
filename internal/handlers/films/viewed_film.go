package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
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
		HandleFilmDetailCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitViewedFilmUserRating:
		parser.ParseFilmUserRating(app, session, requestViewedFilmUserRating, requestViewedFilmReview)
	case states.AwaitViewedFilmReview:
		parser.ParseFilmReview(app, session, requestViewedFilmReview, finishViewedFilmProcess)
	}
}

func finishViewedFilmProcess(app models.App, session *models.Session) {
	HandleUpdateFilm(app, session, HandleFilmDetailCommand)
}

func requestViewedFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestViewedFilmUserRating(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitViewedFilmUserRating)
}

func requestViewedFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestViewedFilmReview(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitViewedFilmReview)
}
