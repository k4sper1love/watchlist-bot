package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleViewedFilmCommand handles the command for marking a film as viewed.
// Sets the viewed status to true and prompts the user to enter their personal rating for the film.
func HandleViewedFilmCommand(app models.App, session *models.Session) {
	session.FilmDetailState.SetViewed(true)
	requestViewedFilmUserRating(app, session)
}

// HandleViewedFilmProcess processes the workflow for marking a film as viewed.
// Handles states like awaiting input for user rating and review.
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

// requestViewedFilmUserRating prompts the user to enter their personal rating for the viewed film.
func requestViewedFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestViewedFilmUserRating(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitViewedFilmUserRating)
}

// requestViewedFilmReview prompts the user to write a review for the viewed film.
func requestViewedFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestViewedFilmReview(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitViewedFilmReview)
}

// finishViewedFilmProcess finalizes the process of marking a film as viewed.
// Calls the Watchlist service to update the film and navigates back to the detailed view.
func finishViewedFilmProcess(app models.App, session *models.Session) {
	HandleUpdateFilm(app, session, HandleFilmDetailCommand)
}
