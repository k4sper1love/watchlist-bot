package films

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
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
		requestKinopoiskToken(app, session)
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
		parseNewFilmFind(app, session)

	case states.ProcessNewFilmAwaitingURL:
		parseNewFilmFromURL(app, session)

	case states.ProcessNewFilmAwaitingTitle:
		parseNewFilmTitle(app, session)

	case states.ProcessNewFilmAwaitingYear:
		parseNewFilmYear(app, session)

	case states.ProcessNewFilmAwaitingGenre:
		parseNewFilmGenre(app, session)

	case states.ProcessNewFilmAwaitingDescription:
		parseNewFilmDescription(app, session)

	case states.ProcessNewFilmAwaitingRating:
		parseNewFilmRating(app, session)

	case states.ProcessNewFilmAwaitingImage:
		parseNewFilmImage(app, session)

	case states.ProcessNewFilmAwaitingComment:
		parseNewFilmComment(app, session)

	case states.ProcessNewFilmAwaitingFilmURL:
		parseNewFilmURL(app, session)

	case states.ProcessNewFilmAwaitingViewed:
		parseNewFilmViewed(app, session)

	case states.ProcessNewFilmAwaitingUserRating:
		parseNewFilmUserRating(app, session)

	case states.ProcessNewFilmAwaitingReview:
		parseNewFilmReview(app, session)

	case states.ProcessNewFilmAwaitingKinopoiskToken:
		parseKinopoiskToken(app, session)
	}
}

func handleNewFilmFind(app models.App, session *models.Session) {
	if session.KinopoiskAPIToken == "" {
		requestKinopoiskToken(app, session)
		return
	}

	app.SendMessage(messages.BuildFilmRequestTitleMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessFindNewFilmAwaitingTitle)
}

func parseNewFilmFind(app models.App, session *models.Session) {
	session.FilmsState.Title = utils.ParseMessageString(app.Update)
	session.FilmsState.CurrentPage = 1

	session.ClearState()
	HandleFindNewFilmCommand(app, session)
}

func handleNewFilmFromURL(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildNewFilmFromURLMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingURL)
}

func parseNewFilmFromURL(app models.App, session *models.Session) {
	url := utils.ParseMessageString(app.Update)
	isKinopoisk := parsing.IsKinopoisk(url)

	if isKinopoisk && session.KinopoiskAPIToken == "" {
		requestKinopoiskToken(app, session)
		return
	}

	film, err := parsing.GetFilmByURL(app, session, url)
	if err != nil {
		handleNewFilmFromURLError(app, session, err, isKinopoisk)
		return
	}

	film.URL = url
	session.FilmDetailState.SetFromFilm(film)
	handleNewFilmUploadImage(app, session)
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

func handleNewFilmUploadImage(app models.App, session *models.Session) {
	imageURL, err := parseAndUploadImageFromURL(app, session.FilmDetailState.ImageURL)
	if err != nil {
		app.SendMessage(messages.BuildImageFailureMessage(session), nil)
	}

	session.FilmDetailState.SetImageURL(imageURL)
	session.FilmsState.CurrentPage = 1
	session.FilmsState.Title = ""

	requestNewFilmComment(app, session)
}

func handleNewFilmManually(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestTitleMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingTitle)
}

func parseNewFilmTitle(app models.App, session *models.Session) {
	if title, ok := parseAndValidateString(app, session, 3, 100); !ok {
		handleNewFilmManually(app, session)
	} else {
		session.FilmDetailState.Title = title
		requestNewFilmYear(app, session)
	}
}

func requestNewFilmYear(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestYearMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingYear)
}

func parseNewFilmYear(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmGenre(app, session)
		return
	}

	if year, ok := parseAndValidateNumber(app, session, 1888, 2100, utils.ParseMessageInt); !ok {
		requestNewFilmYear(app, session)
	} else {
		session.FilmDetailState.Year = year
		requestNewFilmGenre(app, session)
	}
}

func requestNewFilmGenre(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestGenreMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingGenre)
}

func parseNewFilmGenre(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmDescription(app, session)
		return
	}

	if genre, ok := parseAndValidateString(app, session, 0, 100); !ok {
		requestNewFilmGenre(app, session)
	} else {
		session.FilmDetailState.Genre = genre
		requestNewFilmDescription(app, session)
	}
}

func requestNewFilmDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestDescriptionMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingDescription)
}

func parseNewFilmDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmRating(app, session)
		return
	}

	if description, ok := parseAndValidateString(app, session, 0, 1000); !ok {
		requestNewFilmDescription(app, session)
	} else {
		session.FilmDetailState.Description = description
		requestNewFilmRating(app, session)
	}
}

func requestNewFilmRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestRatingMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingRating)
}

func parseNewFilmRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmImage(app, session)
		return
	}

	if rating, ok := parseAndValidateNumber(app, session, 1, 10, utils.ParseMessageFloat); !ok {
		requestNewFilmRating(app, session)
	} else {
		session.FilmDetailState.Rating = rating
		requestNewFilmImage(app, session)
	}
}

func requestNewFilmImage(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestImageMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingImage)
}

func parseNewFilmImage(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmURL(app, session)
		return
	}

	imageURL, err := parseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.BuildImageFailureMessage(session), nil)
		requestNewFilmURL(app, session)
		return
	}

	session.FilmDetailState.SetImageURL(imageURL)
	requestNewFilmURL(app, session)
}

func requestNewFilmURL(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestURLMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingFilmURL)
}

func parseNewFilmURL(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmComment(app, session)
		return
	}

	if url, ok := parseAndValidateURL(app, session); !ok {
		requestNewFilmURL(app, session)
	} else {
		session.FilmDetailState.URL = url
		requestNewFilmComment(app, session)
	}
}

func requestNewFilmComment(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestCommentMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingComment)
}

func parseNewFilmComment(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmViewed(app, session)
		return
	}

	if comment, ok := parseAndValidateString(app, session, 0, 500); !ok {
		requestNewFilmComment(app, session)
	} else {
		session.FilmDetailState.Comment = comment
		requestNewFilmViewed(app, session)
	}
}

func requestNewFilmViewed(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestViewedMessage(session), keyboards.BuildKeyboardWithSurveyAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingViewed)
}

func parseNewFilmViewed(app models.App, session *models.Session) {
	if !utils.IsAgree(app.Update) {
		finishNewFilmProcess(app, session)
		return
	}

	session.FilmDetailState.SetViewed(true)
	requestNewFilmUserRating(app, session)
}

func requestNewFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestUserRatingMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingUserRating)
}

func parseNewFilmUserRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		requestNewFilmReview(app, session)
		return
	}

	if userRating, ok := parseAndValidateNumber(app, session, 1, 10, utils.ParseMessageFloat); !ok {
		requestNewFilmUserRating(app, session)
	} else {
		session.FilmDetailState.UserRating = userRating
		requestNewFilmReview(app, session)
	}
}

func requestNewFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestReviewMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewFilmAwaitingReview)
}

func parseNewFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Update) {
		finishNewFilmProcess(app, session)
		return
	}

	if review, ok := parseAndValidateString(app, session, 0, 500); !ok {
		requestNewFilmReview(app, session)
	} else {
		session.FilmDetailState.Review = review
		finishNewFilmProcess(app, session)
	}
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

func parseAndUploadImageFromURL(app models.App, url string) (string, error) {
	image, err := utils.ParseImageFromURL(url)
	if err != nil {
		return "", err
	}
	return watchlist.UploadImage(app, image)
}

func requestKinopoiskToken(app models.App, session *models.Session) {
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
		return
	}

	app.SendMessage(messages.BuildFilmsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFilmsNew))
}

func parseAndValidateString(app models.App, session *models.Session, min, max int) (string, bool) {
	input := utils.ParseMessageString(app.Update)

	if !utils.IsValidStringLength(input, min, max) {
		validator.HandleInvalidInputLength(app, session, min, max)
		return "", false
	}
	return input, true
}

func parseAndValidateNumber[T int | float64](app models.App, session *models.Session, min T, max T, parser func(*tgbotapi.Update) T) (T, bool) {
	input := parser(app.Update)

	if !utils.IsValidNumberRange(input, min, max) {
		validator.HandleInvalidInputRange(app, session, min, max)
		return 0, false
	}
	return input, true
}

func parseAndValidateURL(app models.App, session *models.Session) (string, bool) {
	input := utils.ParseMessageString(app.Update)

	if !utils.IsValidURL(input) {
		validator.HandleInvalidInputURL(app, session)
		return "", false
	}
	return input, true
}
