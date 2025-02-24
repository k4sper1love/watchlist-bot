package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func parseFilmTitle(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 3, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Title = v })
}

func parseFilmYear(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 1888, 2100, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { s.FilmDetailState.Year = v })
}

func parseFilmGenre(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 0, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Genre = v })
}

func parseFilmDescription(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 0, 1000, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Description = v })
}

func parseFilmRating(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageFloat, utils.IsValidNumberRange[float64], validator.HandleInvalidInputRange[float64], func(s *models.Session, v float64) { s.FilmDetailState.Rating = v })
}

func parseFilmURL(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 5, 2048, utils.ParseMessageString, utils.IsValidURL, validator.HandleInvalidInputURL, func(s *models.Session, v string) { s.FilmDetailState.URL = v })
}

func parseFilmComment(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Comment = v })
}

func parseFilmUserRating(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageFloat, utils.IsValidNumberRange[float64], validator.HandleInvalidInputRange[float64], func(s *models.Session, v float64) { s.FilmDetailState.UserRating = v })
}

func parseFilmReview(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	parser.ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Review = v })
}

func parseFilmImage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := parseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.BuildImageFailureMessage(session), nil)
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

func parseFilmImageWithError(app models.App, session *models.Session, next func(models.App, *models.Session), callback string) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := parseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.BuildImageFailureMessage(session), keyboards.BuildKeyboardWithBack(session, callback))
		session.ClearState()
		return
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

func parseFilmViewed(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.FilmDetailState.SetViewed(utils.IsAgree(app.Update))
	next(app, session)
}

func parseFilmViewedWithFinish(app models.App, session *models.Session, finish, next func(models.App, *models.Session)) {
	if !utils.IsAgree(app.Update) {
		finish(app, session)
		return
	}

	session.FilmDetailState.SetViewed(true)
	next(app, session)
}

func parseFilmFindTitle(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.FilmsState.Title = utils.ParseMessageString(app.Update)
	session.FilmsState.CurrentPage = 1

	session.ClearState()
	next(app, session)
}
