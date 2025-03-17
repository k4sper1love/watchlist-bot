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

// HandleNewFilmCommand handles the command for creating a new film.
// Sends a message with options to create a film manually, from a URL, or by searching.
func HandleNewFilmCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceWay(session), keyboards.FilmNew(session))
}

// HandleNewFilmButtons handles button interactions related to creating a new film.
// Supports actions like going back, creating manually, from a URL, or by searching.
func HandleNewFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallNewFilmBack:
		HandleFilmsCommand(app, session)

	case states.CallNewFilmManually:
		handleNewFilmManually(app, session)

	case states.CallNewFilmFromURL:
		handleNewFilmFromURL(app, session)

	case states.CallNewFilmFind:
		handleNewFilmFind(app, session)

	case states.CallNewFilmChangeKinopoiskToken:
		handleKinopoiskToken(app, session)
	}
}

// HandleNewFilmProcess processes the workflow for creating a new film.
// Handles states like awaiting input for title, year, genre, and other film details.
func HandleNewFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleNewFilmCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitNewFilmFind:
		parser.ParseFindFilmsTitle(app, session, HandleFindNewFilmCommand)

	case states.AwaitNewFilmFromURL:
		parseNewFilmFromURL(app, session)

	case states.AwaitNewFilmKinopoiskToken:
		parser.ParseKinopoiskToken(app, session, HandleNewFilmCommand)

	case states.AwaitNewFilmTitle:
		parser.ParseFilmTitle(app, session, handleNewFilmManually, requestNewFilmYear)

	case states.AwaitNewFilmYear:
		parser.ParseFilmYear(app, session, requestNewFilmYear, requestNewFilmGenre)

	case states.AwaitNewFilmGenre:
		parser.ParseFilmGenre(app, session, requestNewFilmGenre, requestNewFilmDescription)

	case states.AwaitNewFilmDescription:
		parser.ParseFilmDescription(app, session, requestNewFilmDescription, requestNewFilmRating)

	case states.AwaitNewFilmRating:
		parser.ParseFilmRating(app, session, requestNewFilmRating, requestNewFilmImage)

	case states.AwaitNewFilmImage:
		parser.ParseFilmImageFromMessage(app, session, requestNewFilmURL)

	case states.AwaitNewFilmFilmURL:
		parser.ParseFilmURL(app, session, requestNewFilmURL, requestNewFilmComment)

	case states.AwaitNewFilmComment:
		parser.ParseFilmComment(app, session, requestNewFilmComment, requestNewFilmViewed)

	case states.AwaitNewFilmViewed:
		parser.ParseFilmViewedWithFinish(app, session, finishNewFilmProcess, requestNewFilmUserRating)

	case states.AwaitNewFilmUserRating:
		parser.ParseFilmUserRating(app, session, requestNewFilmUserRating, requestNewFilmReview)

	case states.AwaitNewFilmReview:
		parser.ParseFilmReview(app, session, requestNewFilmReview, finishNewFilmProcess)
	}
}

// handleKinopoiskToken prompts the user to enter a Kinopoisk API token.
func handleKinopoiskToken(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestKinopoiskToken(session), keyboards.Cancel(session))
	session.SetState(states.AwaitNewFilmKinopoiskToken)
}

// handleKinopoiskError handles errors related to Kinopoisk API requests.
// Displays appropriate error messages based on the status code.
func handleKinopoiskError(app models.App, session *models.Session, err error) {
	code := client.ParseErrorStatusCode(err)
	if code == 401 || code == 403 {
		app.SendMessage(messages.KinopoiskFailureCode(session, code), keyboards.NewFilmChangeToken(session))
	} else {
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallFilmsNew))
	}
}

// handleNewFilmFind prompts the user to search for a new film by title using Kinopoisk.
func handleNewFilmFind(app models.App, session *models.Session) {
	if session.KinopoiskAPIToken == "" {
		handleKinopoiskToken(app, session)
		return
	}

	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.AwaitNewFilmFind)
}

// handleNewFilmFromURL prompts the user to create a new film by providing a URL.
func handleNewFilmFromURL(app models.App, session *models.Session) {
	app.SendMessage(messages.NewFilmFromURL(session), keyboards.Cancel(session))
	session.SetState(states.AwaitNewFilmFromURL)
}

// parseNewFilmFromURL processes the URL provided by the user to create a new film.
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

// handleNewFilmFromURLError handles errors encountered while processing a film from a URL.
func handleNewFilmFromURLError(app models.App, session *models.Session, err error, isKinopoisk bool) {
	session.ClearState()

	if isKinopoisk {
		handleKinopoiskError(app, session, err)
		return
	}

	app.SendMessage(messages.FilmsFailure(session), nil)
	HandleNewFilmCommand(app, session)
}

// handleNewFilmManually prompts the user to create a new film by entering details manually.
func handleNewFilmManually(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.AwaitNewFilmTitle)
}

// requestNewFilmYear prompts the user to enter the release year of the film.
func requestNewFilmYear(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmYear(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmYear)
}

// requestNewFilmGenre prompts the user to enter the genre of the film.
func requestNewFilmGenre(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmGenre(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmGenre)
}

// requestNewFilmDescription prompts the user to enter a description for the film.
func requestNewFilmDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmDescription(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmDescription)
}

// requestNewFilmRating prompts the user to enter the rating of the film.
func requestNewFilmRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmRating(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmRating)
}

// requestNewFilmImage prompts the user to provide an image for the film.
func requestNewFilmImage(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmImage(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmImage)
}

// requestNewFilmURL prompts the user to provide a URL for the film.
func requestNewFilmURL(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmURL(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmFilmURL)
}

// requestNewFilmComment prompts the user to add a comment for the film.
func requestNewFilmComment(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmComment(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmComment)
}

// requestNewFilmViewed prompts the user to indicate if the film has been viewed.
func requestNewFilmViewed(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmViewed(session), keyboards.SurveyAndCancel(session))
	session.SetState(states.AwaitNewFilmViewed)
}

// requestNewFilmUserRating prompts the user to enter their personal rating for the film.
func requestNewFilmUserRating(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmUserRating(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmUserRating)
}

// requestNewFilmReview prompts the user to write a review for the film.
func requestNewFilmReview(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmReview(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewFilmReview)
}

// finishNewFilmProcess finalizes the creation of a new film.
// Calls the Watchlist service to save the film and navigates to the detailed view.
func finishNewFilmProcess(app models.App, session *models.Session) {
	film, err := createNewFilm(app, session)
	session.ClearAllStates()
	if err != nil {
		app.SendMessage(messages.CreateFilmFailure(session), keyboards.Back(session, states.CallFilmsNew))
		return
	}

	session.FilmDetailState.UpdateFilm(*film)
	HandleFilmDetailCommand(app, session)
}

// createNewFilm creates a new film based on the current session context (user or collection).
func createNewFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	switch session.Context {
	case states.CtxFilm:
		return createNewUserFilm(app, session)
	case states.CtxCollection:
		return createNewCollectionFilm(app, session)
	default:
		return nil, fmt.Errorf("unsupported session context: %v", session.Context)
	}
}

// createNewUserFilm creates a new film associated with the current user using the Watchlist service.
func createNewUserFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	film, err := watchlist.CreateFilm(app, session)
	if err != nil {
		return nil, err
	}

	app.SendMessage(messages.CreateFilmSuccess(session), nil)
	return film, nil
}

// createNewCollectionFilm creates a new film and associates it with a collection using the Watchlist service.
func createNewCollectionFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	collectionFilm, err := watchlist.CreateCollectionFilm(app, session)
	if err != nil {
		return nil, err
	}

	app.SendMessage(messages.CreateCollectionFilmSuccess(session, collectionFilm.Collection.Name), nil)
	return &collectionFilm.Film, nil
}
