package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
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
		parser.ParseFilmURL(app, session, handleUpdateFilmURL, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingImage:
		parser.ParseFilmImageWithError(app, session, finishUpdateFilmProcess, states.CallbackManageFilmSelectUpdate)

	case states.ProcessUpdateFilmAwaitingTitle:
		parser.ParseFilmTitle(app, session, handleUpdateFilmTitle, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingDescription:
		parser.ParseFilmDescription(app, session, handleUpdateFilmDescription, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingGenre:
		parser.ParseFilmGenre(app, session, handleUpdateFilmGenre, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingRating:
		parser.ParseFilmRating(app, session, handleUpdateFilmRating, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingYear:
		parser.ParseFilmYear(app, session, handleUpdateFilmYear, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingComment:
		parser.ParseFilmComment(app, session, handleUpdateFilmComment, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingViewed:
		parser.ParseFilmViewed(app, session, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingUserRating:
		parser.ParseFilmUserRating(app, session, handleUpdateFilmUserRating, finishUpdateFilmProcess)

	case states.ProcessUpdateFilmAwaitingReview:
		parser.ParseFilmReview(app, session, handleUpdateFilmReview, finishUpdateFilmProcess)
	}
}

func finishUpdateFilmProcess(app models.App, session *models.Session) {
	HandleUpdateFilm(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmURL(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestURLMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingURL)
}

func handleUpdateFilmImage(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestImageMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingImage)
}

func handleUpdateFilmTitle(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestTitleMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingTitle)
}

func handleUpdateFilmDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestDescriptionMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingDescription)
}

func handleUpdateFilmGenre(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestGenreMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingGenre)
}

func handleUpdateFilmRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestRatingMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingRating)
}

func handleUpdateFilmYear(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestYearMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingYear)
}

func handleUpdateFilmComment(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestCommentMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingComment)
}

func handleUpdateFilmViewed(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestViewedMessage(session), keyboards.BuildKeyboardWithSurveyAndCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingViewed)
}

func handleUpdateFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestUserRatingMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingUserRating)
}

func handleUpdateFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestReviewMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateFilmAwaitingReview)
}

func HandleUpdateFilm(app models.App, session *models.Session, backFunc func(models.App, *models.Session)) {
	session.FilmDetailState.SyncValues()

	if err := updateFilmAndState(app, session); err != nil {
		app.SendMessage(messages.BuildUpdateFilmFailureMessage(session), nil)
	} else {
		app.SendMessage(messages.BuildUpdateFilmSuccessMessage(session), nil)
	}

	session.ClearAllStates()
	backFunc(app, session)
}

func updateFilmAndState(app models.App, session *models.Session) error {
	film, err := watchlist.UpdateFilm(app, session)
	if err != nil {
		return err
	}

	session.FilmDetailState.UpdateFilmState(*film)
	return nil
}
