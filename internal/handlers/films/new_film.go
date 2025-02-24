package films

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/services/parsing"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleNewFilmCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildChoiceWayMessage(session), keyboards.BuildFilmNewKeyboard(session))
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
		parseFilmFindTitle(app, session, HandleFindNewFilmCommand)

	case states.ProcessNewFilmAwaitingURL:
		parseNewFilmFromURL(app, session)

	case states.ProcessNewFilmAwaitingKinopoiskToken:
		parseKinopoiskToken(app, session)

	case states.ProcessNewFilmAwaitingTitle:
		parseFilmTitle(app, session, handleNewFilmManually, requestNewFilmYear)

	case states.ProcessNewFilmAwaitingYear:
		parseFilmYear(app, session, requestNewFilmYear, requestNewFilmGenre)

	case states.ProcessNewFilmAwaitingGenre:
		parseFilmGenre(app, session, requestNewFilmGenre, requestNewFilmDescription)

	case states.ProcessNewFilmAwaitingDescription:
		parseFilmDescription(app, session, requestNewFilmDescription, requestNewFilmRating)

	case states.ProcessNewFilmAwaitingRating:
		parseFilmRating(app, session, requestNewFilmRating, requestNewFilmImage)

	case states.ProcessNewFilmAwaitingImage:
		parseFilmImage(app, session, requestNewFilmURL)

	case states.ProcessNewFilmAwaitingFilmURL:
		parseFilmURL(app, session, requestNewFilmURL, requestNewFilmComment)

	case states.ProcessNewFilmAwaitingComment:
		parseFilmComment(app, session, requestNewFilmComment, requestNewFilmViewed)

	case states.ProcessNewFilmAwaitingViewed:
		parseFilmViewedWithFinish(app, session, finishNewFilmProcess, requestNewFilmUserRating)

	case states.ProcessNewFilmAwaitingUserRating:
		parseFilmUserRating(app, session, requestNewFilmUserRating, requestNewFilmReview)

	case states.ProcessNewFilmAwaitingReview:
		parseFilmReview(app, session, requestNewFilmReview, finishNewFilmProcess)
	}
}

func handleNewFilmFind(app models.App, session *models.Session) {
	if session.KinopoiskAPIToken == "" {
		handleKinopoiskToken(app, session)
		return
	}

	app.SendMessage(messages.BuildFilmRequestTitleMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessFindNewFilmAwaitingTitle)
}

func handleNewFilmFromURL(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildNewFilmFromURLMessage(session), keyboards.BuildKeyboardWithCancel(session))
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
	parseFilmImage(app, session, requestNewFilmComment)
}

func handleNewFilmFromURLError(app models.App, session *models.Session, err error, isKinopoisk bool) {
	session.ClearState()

	if isKinopoisk {
		handleKinopoiskError(app, session, err)
		return
	}

	app.SendMessage(messages.BuildFilmsFailureMessage(session), nil)
	HandleNewFilmCommand(app, session)
}

func handleNewFilmManually(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestTitleMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingTitle)
}

func requestNewFilmYear(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestYearMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingYear)
}

func requestNewFilmGenre(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestGenreMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingGenre)
}

func requestNewFilmDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestDescriptionMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingDescription)
}

func requestNewFilmRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestRatingMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingRating)
}

func requestNewFilmImage(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestImageMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingImage)
}

func requestNewFilmURL(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestURLMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingFilmURL)
}

func requestNewFilmComment(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestCommentMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingComment)
}

func requestNewFilmViewed(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestViewedMessage(session), keyboards.BuildKeyboardWithSurveyAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingViewed)
}

func requestNewFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestUserRatingMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingUserRating)
}

func requestNewFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestReviewMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingReview)
}

func finishNewFilmProcess(app models.App, session *models.Session) {
	film, err := createNewFilm(app, session)
	session.ClearAllStates()
	if err != nil {
		app.SendMessage(messages.BuildCreateFilmFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFilmsNew))
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

	app.SendMessage(messages.BuildCreateFilmSuccessMessage(session), nil)
	return film, nil
}

func createNewCollectionFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	collectionFilm, err := watchlist.CreateCollectionFilm(app, session)
	if err != nil {
		return nil, err
	}

	app.SendMessage(messages.BuildCreateCollectionFilmSuccessMessage(session, collectionFilm.Collection.Name), nil)
	return &collectionFilm.Film, nil
}

func parseAndUploadImageFromMessage(app models.App) (string, error) {
	image, err := utils.ParseImageFromMessage(app.Bot, app.Update)
	if err != nil {
		return "", err
	}
	return watchlist.UploadImage(app, image)
}

func handleKinopoiskToken(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildKinopoiskTokenMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingKinopoiskToken)
}

func parseKinopoiskToken(app models.App, session *models.Session) {
	session.KinopoiskAPIToken = utils.ParseMessageString(app.Update)
	app.SendMessage(messages.BuildKinopoiskTokenSuccessMessage(session), nil)
	HandleNewFilmCommand(app, session)
}

func handleKinopoiskError(app models.App, session *models.Session, err error) {
	code := client.ParseErrorStatusCode(err)
	if code == 401 || code == 403 {
		app.SendMessage(messages.BuildTokenCodeMessage(session, code), keyboards.BuildNewFilmChangeTokenKeyboard(session))
	} else {
		app.SendMessage(messages.BuildFilmsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFilmsNew))
	}
}
