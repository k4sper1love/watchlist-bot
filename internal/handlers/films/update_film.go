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
		messages.UpdateFilm(session),
		keyboards.FilmUpdate(session),
	)
}

func HandleUpdateFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallUpdateFilmBack:
		HandleFilmDetailCommand(app, session)

	case states.CallUpdateFilmURL:
		handleUpdateFilmURL(app, session)

	case states.CallUpdateFilmImage:
		handleUpdateFilmImage(app, session)

	case states.CallUpdateFilmTitle:
		handleUpdateFilmTitle(app, session)

	case states.CallUpdateFilmDescription:
		handleUpdateFilmDescription(app, session)

	case states.CallUpdateFilmGenre:
		handleUpdateFilmGenre(app, session)

	case states.CallUpdateFilmRating:
		handleUpdateFilmRating(app, session)

	case states.CallUpdateFilmYear:
		handleUpdateFilmYear(app, session)

	case states.CallUpdateFilmComment:
		handleUpdateFilmComment(app, session)

	case states.CallUpdateFilmViewed:
		handleUpdateFilmViewed(app, session)

	case states.CallUpdateFilmUserRating:
		handleUpdateFilmUserRating(app, session)

	case states.CallUpdateFilmReview:
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
	case states.AwaitUpdateFilmURL:
		parser.ParseFilmURL(app, session, handleUpdateFilmURL, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmImage:
		parser.ParseFilmImageFromMessageWithError(app, session, finishUpdateFilmProcess, states.CallManageFilmUpdate)

	case states.AwaitUpdateFilmTitle:
		parser.ParseFilmTitle(app, session, handleUpdateFilmTitle, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmDescription:
		parser.ParseFilmDescription(app, session, handleUpdateFilmDescription, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmGenre:
		parser.ParseFilmGenre(app, session, handleUpdateFilmGenre, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmRating:
		parser.ParseFilmRating(app, session, handleUpdateFilmRating, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmYear:
		parser.ParseFilmYear(app, session, handleUpdateFilmYear, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmComment:
		parser.ParseFilmComment(app, session, handleUpdateFilmComment, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmViewed:
		parser.ParseFilmViewed(app, session, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmUserRating:
		parser.ParseFilmUserRating(app, session, handleUpdateFilmUserRating, finishUpdateFilmProcess)

	case states.AwaitUpdateFilmReview:
		parser.ParseFilmReview(app, session, handleUpdateFilmReview, finishUpdateFilmProcess)
	}
}

func finishUpdateFilmProcess(app models.App, session *models.Session) {
	HandleUpdateFilm(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmURL(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmURL(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmURL)
}

func handleUpdateFilmImage(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmImage(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmImage)
}

func handleUpdateFilmTitle(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmTitle)
}

func handleUpdateFilmDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmDescription(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmDescription)
}

func handleUpdateFilmGenre(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmGenre(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmGenre)
}

func handleUpdateFilmRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmRating(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmRating)
}

func handleUpdateFilmYear(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmYear(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmYear)
}

func handleUpdateFilmComment(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmComment(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmComment)
}

func handleUpdateFilmViewed(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmViewed(session), keyboards.SurveyAndCancel(session))
	session.SetState(states.AwaitUpdateFilmViewed)
}

func handleUpdateFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmUserRating(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmUserRating)
}

func handleUpdateFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmReview(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateFilmReview)
}

func HandleUpdateFilm(app models.App, session *models.Session, backFunc func(models.App, *models.Session)) {
	session.FilmDetailState.SyncValues()

	if err := updateFilmAndState(app, session); err != nil {
		app.SendMessage(messages.UpdateFilmFailure(session), nil)
	} else {
		app.SendMessage(messages.UpdateFilmSuccess(session), nil)
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
