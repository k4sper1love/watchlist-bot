package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func ParseFilmTitle(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 3, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Title = v })
}

func ParseFilmYear(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1888, 2100, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { s.FilmDetailState.Year = v })
}

func ParseFilmGenre(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Genre = v })
}

func ParseFilmDescription(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 1000, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Description = v })
}

func ParseFilmRating(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageFloat, utils.IsValidNumberRange[float64], validator.HandleInvalidInputRange[float64], func(s *models.Session, v float64) { s.FilmDetailState.Rating = v })
}

func ParseFilmURL(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 5, 2048, utils.ParseMessageString, utils.IsValidURL, validator.HandleInvalidInputURL, func(s *models.Session, v string) { s.FilmDetailState.URL = v })
}

func ParseFilmComment(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Comment = v })
}

func ParseFilmUserRating(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageFloat, utils.IsValidNumberRange[float64], validator.HandleInvalidInputRange[float64], func(s *models.Session, v float64) { s.FilmDetailState.UserRating = v })
}

func ParseFilmReview(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Review = v })
}

func ParseFilmImageFromMessage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := ParseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), nil)
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

func ParseFilmImageFromMessageWithError(app models.App, session *models.Session, next func(models.App, *models.Session), callback string) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := ParseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), keyboards.BuildKeyboardWithBack(session, callback))
		session.ClearState()
		return
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

func ParseFilmImageFromURL(app models.App, session *models.Session, imageURL string, next func(models.App, *models.Session)) {
	imageURL, err := ParseAndUploadImageFromURL(app, imageURL)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), nil)
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

func ParseFilmViewed(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.FilmDetailState.SetViewed(utils.IsAgree(app.Update))
	next(app, session)
}

func ParseFilmViewedWithFinish(app models.App, session *models.Session, finish, next func(models.App, *models.Session)) {
	if !utils.IsAgree(app.Update) {
		finish(app, session)
		return
	}

	session.FilmDetailState.SetViewed(true)
	next(app, session)
}

func ParseFilmFindTitle(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.FilmsState.Title = utils.ParseMessageString(app.Update)
	session.FilmsState.CurrentPage = 1

	session.ClearState()
	next(app, session)
}
