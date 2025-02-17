package films

import (
	"fmt"
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
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleNewFilmCommand(app models.App, session *models.Session) {
	choiceMsg := translator.Translate(session.Lang, "choiceWay", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>", choiceMsg)

	keyboard := keyboards.BuildFilmNewKeyboard(session)
	app.SendMessage(msg, keyboard)
}

func HandleNewFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
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
	if utils.IsCancel(app.Upd) {
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

	msg := "❓" + translator.Translate(session.Lang, "filmRequestTitle", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessFindNewFilmAwaitingTitle)
}

func parseNewFilmFind(app models.App, session *models.Session) {
	title := utils.ParseMessageString(app.Upd)

	session.FilmsState.Title = title
	session.FilmsState.CurrentPage = 1

	session.ClearState()

	HandleFindNewFilmCommand(app, session)
}

func handleNewFilmFromURL(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	part1 := translator.Translate(session.Lang, "filmRequestLink", nil, nil)
	part2 := translator.Translate(session.Lang, "supportedServices", nil, nil)
	supportedServices := parsing.GetSupportedServicesInline()

	msg := fmt.Sprintf("❓<b>%s</b>\n\n%s:\n<i>%s</i>", part1, part2, supportedServices)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessNewFilmAwaitingURL)
}

func parseNewFilmFromURL(app models.App, session *models.Session) {
	url := utils.ParseMessageString(app.Upd)
	isKinopoisk := parsing.IsKinopoisk(url)

	if isKinopoisk && session.KinopoiskAPIToken == "" {
		requestKinopoiskToken(app, session)
		return
	}

	film, err := parsing.GetFilmByURL(app, session, url)
	if err != nil {
		if isKinopoisk {
			handleKinopoiskError(app, session, err)
			return
		}
		msg := "🚨 " + translator.Translate(session.Lang, "getFilmFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleNewFilmCommand(app, session)
		return
	}

	film.URL = url
	session.FilmDetailState.SetFromFilm(film)

	imageURL, err := parseAndUploadImageFromURL(app, film.ImageURL)
	if err != nil {
		msg := "⚠️ " + translator.Translate(session.Lang, "getImageFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.FilmDetailState.SetImageURL("")
		requestNewFilmComment(app, session)
		return
	}
	session.FilmDetailState.SetImageURL(imageURL)

	requestNewFilmComment(app, session)
}

func handleNewFilmManually(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestTitle", nil, nil)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessNewFilmAwaitingTitle)
}

func parseNewFilmTitle(app models.App, session *models.Session) {
	title := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(title, 3, 100); !ok {
		validator.HandleInvalidInputLength(app, session, 3, 100)
		handleNewFilmManually(app, session)
		return
	}
	session.FilmDetailState.Title = utils.ParseMessageString(app.Upd)

	requestNewFilmYear(app, session)
}

func requestNewFilmYear(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestYear", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingYear)
}

func parseNewFilmYear(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Year = 0
		requestNewFilmGenre(app, session)
		return
	}

	year := utils.ParseMessageInt(app.Upd)
	if ok := utils.ValidNumberRange(year, 1888, 2100); !ok {
		validator.HandleInvalidInputRange(app, session, 1888, 2100)
		requestNewFilmYear(app, session)
		return
	}
	session.FilmDetailState.Year = year

	requestNewFilmGenre(app, session)
}

func requestNewFilmGenre(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestGenre", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingGenre)
}

func parseNewFilmGenre(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Genre = ""
		requestNewFilmDescription(app, session)
		return
	}

	genre := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(genre, 0, 100); !ok {
		validator.HandleInvalidInputLength(app, session, 0, 100)
		requestNewFilmGenre(app, session)
		return
	}
	session.FilmDetailState.Genre = genre

	requestNewFilmDescription(app, session)
}

func requestNewFilmDescription(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestDescription", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingDescription)
}

func parseNewFilmDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Description = ""
		requestNewFilmRating(app, session)
		return
	}

	description := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(description, 0, 1000); !ok {
		validator.HandleInvalidInputLength(app, session, 0, 1000)
		requestNewFilmDescription(app, session)
		return
	}
	session.FilmDetailState.Description = description

	requestNewFilmRating(app, session)
}

func requestNewFilmRating(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestRating", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingRating)
}

func parseNewFilmRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Rating = 0
		requestNewFilmImage(app, session)
		return
	}

	rating := utils.ParseMessageFloat(app.Upd)
	if ok := utils.ValidNumberRange(rating, 1, 10); !ok {
		validator.HandleInvalidInputRange(app, session, 1, 10)
		requestNewFilmRating(app, session)
		return
	}
	session.FilmDetailState.Rating = rating

	requestNewFilmImage(app, session)
}

func requestNewFilmImage(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestImage", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingImage)
}

func parseNewFilmImage(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		requestNewFilmURL(app, session)
		return
	}

	imageURL, err := parseAndUploadImageFromMessage(app)
	if err != nil {
		msg := "⚠️ " + translator.Translate(session.Lang, "getImageFailure", nil, nil)
		app.SendMessage(msg, nil)
		requestNewFilmURL(app, session)
		return
	}

	session.FilmDetailState.SetImageURL(imageURL)

	requestNewFilmURL(app, session)
}

func requestNewFilmURL(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestLink", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingFilmURL)
}

func parseNewFilmURL(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.URL = ""
		requestNewFilmComment(app, session)
		return
	}

	u := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidURL(u); !ok {
		validator.HandleInvalidInputURL(app, session)
		requestNewFilmURL(app, session)
		return
	}
	session.FilmDetailState.URL = u

	requestNewFilmComment(app, session)
}

func requestNewFilmComment(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestComment", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingComment)
}

func parseNewFilmComment(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Comment = ""
		requestNewFilmViewed(app, session)
		return
	}

	comment := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(comment, 0, 500); !ok {
		validator.HandleInvalidInputLength(app, session, 0, 500)
		requestNewFilmComment(app, session)
		return
	}
	session.FilmDetailState.Comment = comment

	requestNewFilmViewed(app, session)
}

func requestNewFilmViewed(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSurvey().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestViewed", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingViewed)
}

func parseNewFilmViewed(app models.App, session *models.Session) {
	switch utils.IsAgree(app.Upd) {
	case true:
		session.FilmDetailState.IsViewed = true
		requestNewFilmUserRating(app, session)

	case false:
		session.FilmDetailState.IsViewed = false
		session.FilmDetailState.UserRating = 0
		session.FilmDetailState.Review = ""
		finishNewFilmProcess(app, session)
	}
}

func requestNewFilmUserRating(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestUserRating", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingUserRating)
}

func parseNewFilmUserRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.UserRating = 0
		requestNewFilmReview(app, session)
		return
	}

	userRating := utils.ParseMessageFloat(app.Upd)
	if ok := utils.ValidNumberRange(userRating, 1, 10); !ok {
		validator.HandleInvalidInputRange(app, session, 1, 10)
		requestNewFilmUserRating(app, session)
		return
	}
	session.FilmDetailState.UserRating = userRating

	requestNewFilmReview(app, session)
}

func requestNewFilmReview(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestReview", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingReview)
}

func parseNewFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Review = ""
		finishNewFilmProcess(app, session)
		return
	}

	review := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(review, 0, 500); !ok {
		validator.HandleInvalidInputLength(app, session, 0, 500)
		requestNewFilmReview(app, session)
		return
	}
	session.FilmDetailState.Review = review

	finishNewFilmProcess(app, session)
}

func finishNewFilmProcess(app models.App, session *models.Session) {
	film, err := CreateNewFilm(app, session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "createFilmFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackFilmsNew).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		session.ClearAllStates()
		return
	}

	session.ClearAllStates()

	session.FilmDetailState.Film = *film
	session.FilmDetailState.ClearIndex()

	HandleFilmsDetailCommand(app, session)
}

func CreateNewFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
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

	msg := "🎬 " + translator.Translate(session.Lang, "createFilmSuccess", nil, nil)
	app.SendMessage(msg, nil)

	return film, nil
}

func createNewCollectionFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	collectionFilm, err := watchlist.CreateCollectionFilm(app, session)
	if err != nil {
		return nil, err
	}

	msg := "🎬 " + translator.Translate(session.Lang, "createCollectionFilmSuccess", map[string]interface{}{
		"Collection": collectionFilm.Collection.Name,
	}, nil)
	app.SendMessage(msg, nil)

	return &collectionFilm.Film, nil
}

func parseAndUploadImageFromMessage(app models.App) (string, error) {
	image, err := utils.ParseImageFromMessage(app.Bot, app.Upd)
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
	msg := messages.BuildKinopoiskTokenMessage(session)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingKinopoiskToken)
}

func parseKinopoiskToken(app models.App, session *models.Session) {
	token := utils.ParseMessageString(app.Upd)

	session.KinopoiskAPIToken = token

	msg := messages.BuildKinopoiskTokenSuccessMessage(session)
	app.SendMessage(msg, nil)

	HandleNewFilmCommand(app, session)
}

func handleKinopoiskError(app models.App, session *models.Session, err error) {
	session.ClearAllStates()

	msg := "🚨 "
	switch client.ParseErrorStatusCode(err) {
	case 401:
		msg += translator.Translate(session.Lang, "tokenFailure", nil, nil)
	case 403:
		msg += translator.Translate(session.Lang, "tokenLimit", nil, nil)
	default:
		msg += translator.Translate(session.Lang, "getFilmFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackFilmsNew).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	keyboard := keyboards.NewKeyboard().AddChangeToken().AddBack(states.CallbackFilmsNew).Build(session.Lang)

	app.SendMessage(msg, keyboard)
}
