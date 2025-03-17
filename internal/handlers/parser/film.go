package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// ParseFilmTitle processes the input for a film's title.
// Validates the title length and retries if the input is invalid.
// Stores the validated title in the session's FilmDetailState.
func ParseFilmTitle(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 3, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Title = v })
}

// ParseFilmYear processes the input for a film's release year.
// Validates the year range and retries if the input is invalid.
// Stores the validated year in the session's FilmDetailState.
func ParseFilmYear(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1888, 2100, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { s.FilmDetailState.Year = v })
}

// ParseFilmGenre processes the input for a film's genre.
// Validates the genre length and retries if the input is invalid.
// Stores the validated genre in the session's FilmDetailState.
func ParseFilmGenre(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Genre = v })
}

// ParseFilmDescription processes the input for a film's description.
// Validates the description length and retries if the input is invalid.
// Stores the validated description in the session's FilmDetailState.
func ParseFilmDescription(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 1000, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Description = v })
}

// ParseFilmRating processes the input for a film's rating.
// Validates the rating range and retries if the input is invalid.
// Stores the validated rating in the session's FilmDetailState.
func ParseFilmRating(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageFloat, utils.IsValidNumberRange[float64], validator.HandleInvalidInputRange[float64], func(s *models.Session, v float64) { s.FilmDetailState.Rating = v })
}

// ParseFilmURL processes the input for a film's URL.
// Validates the URL format and retries if the input is invalid.
// Stores the validated URL in the session's FilmDetailState.
func ParseFilmURL(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 5, 2048, utils.ParseMessageString, utils.IsValidURL, validator.HandleInvalidInputURL, func(s *models.Session, v string) { s.FilmDetailState.URL = v })
}

// ParseFilmComment processes the input for a film's comment.
// Validates the comment length and retries if the input is invalid.
// Stores the validated comment in the session's FilmDetailState.
func ParseFilmComment(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Comment = v })
}

// ParseFilmUserRating processes the input for a user's personal rating of a film.
// Validates the rating range and retries if the input is invalid.
// Stores the validated user rating in the session's FilmDetailState.
func ParseFilmUserRating(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageFloat, utils.IsValidNumberRange[float64], validator.HandleInvalidInputRange[float64], func(s *models.Session, v float64) { s.FilmDetailState.UserRating = v })
}

// ParseFilmReview processes the input for a film's review.
// Validates the review length and retries if the input is invalid.
// Stores the validated review in the session's FilmDetailState.
func ParseFilmReview(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.FilmDetailState.Review = v })
}

// ParseFilmImageFromMessage processes the input for uploading a film's image from a message.
// Skips the step if the user chooses to skip; otherwise, uploads the image and stores its URL.
func ParseFilmImageFromMessage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := UploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), nil)
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

// ParseFilmImageFromMessageWithError processes the input for uploading a film's image from a message with error handling.
// Skips the step if the user chooses to skip; otherwise, uploads the image and handles errors with a callback.
func ParseFilmImageFromMessageWithError(app models.App, session *models.Session, next func(models.App, *models.Session), callback string) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := UploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), keyboards.Back(session, callback))
		session.ClearState()
		return
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

// ParseFilmImageFromURL processes the input for uploading a film's image from a URL.
// Uploads the image and stores its URL in the session's FilmDetailState.
func ParseFilmImageFromURL(app models.App, session *models.Session, imageURL string, next func(models.App, *models.Session)) {
	imageURL, err := UploadImageFromURL(app, imageURL)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), nil)
	}

	session.FilmDetailState.SetImageURL(imageURL)
	next(app, session)
}

// ParseFilmViewed processes the input for marking a film as viewed.
// Sets the viewed status in the session's FilmDetailState based on the user's choice.
func ParseFilmViewed(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.FilmDetailState.SetViewed(utils.IsAgree(app.Update))
	next(app, session)
}

// ParseFilmViewedWithFinish processes the input for marking a film as viewed with an optional finish step.
// If the user agrees, marks the film as viewed and proceeds; otherwise, calls the finish function.
func ParseFilmViewedWithFinish(app models.App, session *models.Session, finish, next func(models.App, *models.Session)) {
	if !utils.IsAgree(app.Update) {
		finish(app, session)
		return
	}

	session.FilmDetailState.SetViewed(true)
	next(app, session)
}

// ParseFindFilmsTitle processes the input for searching films by title.
// Parses the input string, resets the current page, and clears the session state before proceeding.
func ParseFindFilmsTitle(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.FilmsState.Title = utils.ParseMessageString(app.Update)
	session.FilmsState.CurrentPage = 1

	session.ClearState()
	next(app, session)
}
