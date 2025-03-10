package films

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/services/parsing"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleNewFilmCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceWay(session), keyboards.FilmNew(session))
}

func HandleNewFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackNewFilmSelectBack:
		HandleFilmsCommand(app, session)

	case states.CallbackNewFilmSelectManually:
		handleNewFilmManually(app, session)

	case states.CallbackNewFilmSelectFromURL:
		handleNewFilmFromURL(app, session)

	case states.CallbackNewFilmSelectFind:
		handleNewFilmFind(app, session)

	case states.CallbackNewFilmSelectChangeKinopoiskToken:
		handleKinopoiskToken(app, session)
	}
}

func HandleNewFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleNewFilmCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFindNewFilmAwaitingTitle:
		parser.ParseFilmFindTitle(app, session, HandleFindNewFilmCommand)

	case states.ProcessNewFilmAwaitingURL:
		parseNewFilmFromURL(app, session)

	case states.ProcessNewFilmAwaitingKinopoiskToken:
		parseKinopoiskToken(app, session)

	case states.ProcessNewFilmAwaitingTitle:
		parser.ParseFilmTitle(app, session, handleNewFilmManually, requestNewFilmYear)

	case states.ProcessNewFilmAwaitingYear:
		parser.ParseFilmYear(app, session, requestNewFilmYear, requestNewFilmGenre)

	case states.ProcessNewFilmAwaitingGenre:
		parser.ParseFilmGenre(app, session, requestNewFilmGenre, requestNewFilmDescription)

	case states.ProcessNewFilmAwaitingDescription:
		parser.ParseFilmDescription(app, session, requestNewFilmDescription, requestNewFilmRating)

	case states.ProcessNewFilmAwaitingRating:
		parser.ParseFilmRating(app, session, requestNewFilmRating, requestNewFilmImage)

	case states.ProcessNewFilmAwaitingImage:
		parser.ParseFilmImageFromMessage(app, session, requestNewFilmURL)

	case states.ProcessNewFilmAwaitingFilmURL:
		parser.ParseFilmURL(app, session, requestNewFilmURL, requestNewFilmComment)

	case states.ProcessNewFilmAwaitingComment:
		parser.ParseFilmComment(app, session, requestNewFilmComment, requestNewFilmViewed)

	case states.ProcessNewFilmAwaitingViewed:
		parser.ParseFilmViewedWithFinish(app, session, finishNewFilmProcess, requestNewFilmUserRating)

	case states.ProcessNewFilmAwaitingUserRating:
		parser.ParseFilmUserRating(app, session, requestNewFilmUserRating, requestNewFilmReview)

	case states.ProcessNewFilmAwaitingReview:
		parser.ParseFilmReview(app, session, requestNewFilmReview, finishNewFilmProcess)
	}
}

func handleKinopoiskToken(app models.App, session *models.Session) {
	app.SendMessage(messages.KinopoiskToken(session), keyboards.Cancel(session))
	session.SetState(states.ProcessNewFilmAwaitingKinopoiskToken)
}

func parseKinopoiskToken(app models.App, session *models.Session) {
	session.KinopoiskAPIToken = utils.ParseMessageString(app.Update)
	app.SendMessage(messages.KinopoiskTokenSuccess(session), nil)
	HandleNewFilmCommand(app, session)
}

func handleKinopoiskError(app models.App, session *models.Session, err error) {
	code := client.ParseErrorStatusCode(err)
	if code == 401 || code == 403 {
		app.SendMessage(messages.KinopoiskFailureCode(session, code), keyboards.NewFilmChangeToken(session))
	} else {
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallbackFilmsNew))
	}
}

func handleNewFilmFind(app models.App, session *models.Session) {
	if session.KinopoiskAPIToken == "" {
		handleKinopoiskToken(app, session)
		return
	}

	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.ProcessFindNewFilmAwaitingTitle)
}

func handleNewFilmFromURL(app models.App, session *models.Session) {
	app.SendMessage(messages.NewFilmFromURL(session), keyboards.Cancel(session))
	session.SetState(states.ProcessNewFilmAwaitingURL)
}

func parseNewFilmFromURL(app models.App, session *models.Session) {
	url := utils.ParseMessageString(app.Update)
	isKinopoisk := parsing.IsKinopoisk(url)

	if isKinopoisk && session.KinopoiskAPIToken == "" {
		handleKinopoiskToken(app, session)
		return
	}

	film, err := parsing.GetFilmByURL(app, session, url)
	if err != nil {
		handleNewFilmFromURLError(app, session, err, isKinopoisk)
		return
	}

	film.URL = url
	session.FilmDetailState.SetFromFilm(film)
	parser.ParseFilmImageFromMessage(app, session, requestNewFilmComment)
}

func handleNewFilmFromURLError(app models.App, session *models.Session, err error, isKinopoisk bool) {
	session.ClearState()

	if isKinopoisk {
		handleKinopoiskError(app, session, err)
		return
	}

	app.SendMessage(messages.FilmsFailure(session), nil)
	HandleNewFilmCommand(app, session)
}

func handleNewFilmManually(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.ProcessNewFilmAwaitingTitle)
}

func requestNewFilmYear(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmYear(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingYear)
}

func requestNewFilmGenre(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmGenre(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingGenre)
}

func requestNewFilmDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmDescription(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingDescription)
}

func requestNewFilmRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmRating(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingRating)
}

func requestNewFilmImage(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmImage(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingImage)
}

func requestNewFilmURL(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmURL(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingFilmURL)
}

func requestNewFilmComment(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmComment(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingComment)
}

func requestNewFilmViewed(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmViewed(session), keyboards.SurveyAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingViewed)
}

func requestNewFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmUserRating(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingUserRating)
}

func requestNewFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmReview(session), keyboards.SkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingReview)
}

func finishNewFilmProcess(app models.App, session *models.Session) {
	film, err := createNewFilm(app, session)
	session.ClearAllStates()
	if err != nil {
		app.SendMessage(messages.CreateFilmFailure(session), keyboards.Back(session, states.CallbackFilmsNew))
		return
	}

	session.FilmDetailState.UpdateFilmState(*film)
	HandleFilmsDetailCommand(app, session)
}

func createNewFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	switch session.Context {
	case states.ContextFilm:
		return createNewUserFilm(app, session)
	case states.ContextCollection:
		return createNewCollectionFilm(app, session)
	default:
		return nil, fmt.Errorf("unsupported session context: %v", session.Context)
	}
}

func createNewUserFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	film, err := watchlist.CreateFilm(app, session)
	if err != nil {
		return nil, err
	}

	app.SendMessage(messages.CreateFilmSuccess(session), nil)
	return film, nil
}

func createNewCollectionFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	collectionFilm, err := watchlist.CreateCollectionFilm(app, session)
	if err != nil {
		return nil, err
	}

	app.SendMessage(messages.CreateCollectionFilmSuccess(session, collectionFilm.Collection.Name), nil)
	return &collectionFilm.Film, nil
}
