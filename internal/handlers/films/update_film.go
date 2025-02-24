package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleUpdateFilmCommand(app models.App, session *models.Session) {
	app.SendImage(
		session.FilmDetailState.Film.ImageURL,
		messages.BuildUpdateFilmMessage(session),
		keyboards.BuildFilmUpdateKeyboard(session),
	)
}

func HandleUpdateFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackUpdateFilmSelectBack:
		HandleFilmsDetailCommand(app, session)

	case states.CallbackUpdateFilmSelectURL:
		handleUpdateFilmURL(app, session)

	case states.CallbackUpdateFilmSelectImage:
		handleUpdateFilmImage(app, session)

	case states.CallbackUpdateFilmSelectTitle:
		handleUpdateFilmTitle(app, session)

	case states.CallbackUpdateFilmSelectDescription:
		handleUpdateFilmDescription(app, session)

	case states.CallbackUpdateFilmSelectGenre:
		handleUpdateFilmGenre(app, session)

	case states.CallbackUpdateFilmSelectRating:
		handleUpdateFilmRating(app, session)

	case states.CallbackUpdateFilmSelectYear:
		handleUpdateFilmYear(app, session)

	case states.CallbackUpdateFilmSelectComment:
		handleUpdateFilmComment(app, session)

	case states.CallbackUpdateFilmSelectViewed:
		handleUpdateFilmViewed(app, session)

	case states.CallbackUpdateFilmSelectUserRating:
		handleUpdateFilmUserRating(app, session)

	case states.CallbackUpdateFilmSelectReview:
		handleUpdateFilmReview(app, session)
	}
}

func HandleUpdateFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleUpdateFilmCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessUpdateFilmAwaitingURL:
		parseUpdateFilmURL(app, session)

	case states.ProcessUpdateFilmAwaitingImage:
		parseUpdateFilmImage(app, session)

	case states.ProcessUpdateFilmAwaitingTitle:
		parseUpdateFilmTitle(app, session)

	case states.ProcessUpdateFilmAwaitingDescription:
		parseUpdateFilmDescription(app, session)

	case states.ProcessUpdateFilmAwaitingGenre:
		parseUpdateFilmGenre(app, session)

	case states.ProcessUpdateFilmAwaitingRating:
		parseUpdateFilmRating(app, session)

	case states.ProcessUpdateFilmAwaitingYear:
		parseUpdateFilmYear(app, session)

	case states.ProcessUpdateFilmAwaitingComment:
		parseUpdateFilmComment(app, session)

	case states.ProcessUpdateFilmAwaitingViewed:
		parseUpdateFilmViewed(app, session)

	case states.ProcessUpdateFilmAwaitingUserRating:
		parseUpdateFilmUserRating(app, session)

	case states.ProcessUpdateFilmAwaitingReview:
		parseUpdateFilmReview(app, session)
	}
}

func handleUpdateFilmURL(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestURLMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingURL)
}

func parseUpdateFilmURL(app models.App, session *models.Session) {
	if url, ok := parseAndValidateURL(app, session); !ok {
		handleUpdateFilmURL(app, session)
	} else {
		session.FilmDetailState.URL = url
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmImage(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestImageMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingImage)
}

func parseUpdateFilmImage(app models.App, session *models.Session) {
	imageURL, err := parseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.BuildImageFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackManageFilmSelectUpdate))
		session.ClearState()
		return
	}

	session.FilmDetailState.SetImageURL(imageURL)
	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmTitle(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestTitleMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingTitle)
}

func parseUpdateFilmTitle(app models.App, session *models.Session) {
	if title, ok := parseAndValidateString(app, session, 3, 100); !ok {
		handleUpdateFilmTitle(app, session)
	} else {
		session.FilmDetailState.Title = title
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestDescriptionMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingDescription)
}

func parseUpdateFilmDescription(app models.App, session *models.Session) {
	if description, ok := parseAndValidateString(app, session, 0, 1000); !ok {
		handleUpdateFilmDescription(app, session)
	} else {
		session.FilmDetailState.Description = description
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmGenre(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestGenreMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingGenre)
}

func parseUpdateFilmGenre(app models.App, session *models.Session) {
	if genre, ok := parseAndValidateString(app, session, 0, 100); !ok {
		handleUpdateFilmGenre(app, session)
	} else {
		session.FilmDetailState.Genre = genre
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestRatingMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingRating)
}

func parseUpdateFilmRating(app models.App, session *models.Session) {
	if rating, ok := parseAndValidateNumber(app, session, 1, 10, utils.ParseMessageFloat); !ok {
		handleUpdateFilmRating(app, session)
	} else {
		session.FilmDetailState.Rating = rating
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmYear(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestYearMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingYear)
}

func parseUpdateFilmYear(app models.App, session *models.Session) {
	if year, ok := parseAndValidateNumber(app, session, 1888, 2100, utils.ParseMessageInt); !ok {
		handleUpdateFilmYear(app, session)
	} else {
		session.FilmDetailState.Year = year
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmComment(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestCommentMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingComment)
}

func parseUpdateFilmComment(app models.App, session *models.Session) {
	if comment, ok := parseAndValidateString(app, session, 0, 500); !ok {
		handleUpdateFilmComment(app, session)
	} else {
		session.FilmDetailState.Comment = comment
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmViewed(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestViewedMessage(session), keyboards.BuildKeyboardWithSurveyAndCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingViewed)
}

func parseUpdateFilmViewed(app models.App, session *models.Session) {
	session.FilmDetailState.SetViewed(utils.IsAgree(app.Update))
	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestUserRatingMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingUserRating)
}

func parseUpdateFilmUserRating(app models.App, session *models.Session) {
	if userRating, ok := parseAndValidateNumber(app, session, 1, 10, utils.ParseMessageFloat); !ok {
		handleUpdateFilmUserRating(app, session)
	} else {
		session.FilmDetailState.UserRating = userRating
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func handleUpdateFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestReviewMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingReview)
}

func parseUpdateFilmReview(app models.App, session *models.Session) {
	if review, ok := parseAndValidateString(app, session, 0, 500); !ok {
		handleUpdateFilmReview(app, session)
	} else {
		session.FilmDetailState.Review = review
		finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
	}
}

func updateFilm(app models.App, session *models.Session) error {
	film, err := watchlist.UpdateFilm(app, session)
	if err != nil {
		return err
	}

	session.FilmDetailState.UpdateFilmState(*film)
	return nil
}

func finishUpdateFilmProcess(app models.App, session *models.Session, backFunc func(models.App, *models.Session)) {
	session.FilmDetailState.SyncValues()

	if err := updateFilm(app, session); err != nil {
		app.SendMessage(messages.BuildUpdateFilmFailureMessage(session), nil)
	} else {
		app.SendMessage(messages.BuildUpdateFilmSuccessMessage(session), nil)
	}

	session.ClearAllStates()
	backFunc(app, session)
}
