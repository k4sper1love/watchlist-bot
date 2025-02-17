package films

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
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
	}
}

func handleNewFilmFind(app models.App, session *models.Session) {
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

	film, err := parsing.GetFilmByURL(app, session, url)
	if err != nil {
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
	session.FilmDetailState.Title = utils.ParseMessageString(app.Upd)

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestYear", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingYear)
}

func parseNewFilmYear(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Year = 0
	} else {
		session.FilmDetailState.Year = utils.ParseMessageInt(app.Upd)
	}

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestGenre", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingGenre)
}

func parseNewFilmGenre(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Genre = ""
	} else {
		session.FilmDetailState.Genre = utils.ParseMessageString(app.Upd)
	}

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestDescription", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingDescription)
}

func parseNewFilmDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Description = ""
	} else {
		session.FilmDetailState.Description = utils.ParseMessageString(app.Upd)
	}

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestRating", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingRating)
}

func parseNewFilmRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Rating = 0
	} else {
		session.FilmDetailState.Rating = utils.ParseMessageFloat(app.Upd)
	}

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
	} else {
		session.FilmDetailState.URL = utils.ParseMessageString(app.Upd)
	}

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
	} else {
		session.FilmDetailState.Comment = utils.ParseMessageString(app.Upd)
	}

	keyboard := keyboards.NewKeyboard().AddSurvey().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestViewed", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingViewed)
}

func parseNewFilmViewed(app models.App, session *models.Session) {
	switch utils.IsAgree(app.Upd) {
	case true:
		session.FilmDetailState.IsViewed = true
		keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

		msg := "❓" + translator.Translate(session.Lang, "filmRequestUserRating", nil, nil)

		app.SendMessage(msg, keyboard)
		session.SetState(states.ProcessNewFilmAwaitingUserRating)

	case false:
		session.FilmDetailState.IsViewed = false
		session.FilmDetailState.UserRating = 0
		session.FilmDetailState.Review = ""

		finishNewFilmProcess(app, session)
	}
}

func parseNewFilmUserRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.UserRating = 0
	} else {
		session.FilmDetailState.UserRating = utils.ParseMessageFloat(app.Upd)
	}

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "filmRequestReview", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingReview)
}

func parseNewFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Review = ""
	} else {
		session.FilmDetailState.Review = utils.ParseMessageString(app.Upd)
	}

	finishNewFilmProcess(app, session)
}

func finishNewFilmProcess(app models.App, session *models.Session) {
	film, err := CreateNewFilm(app, session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "createFilmFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleNewFilmCommand(app, session)
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
