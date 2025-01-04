package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleUpdateFilmCommand(app models.App, session *models.Session) {
	film := session.FilmDetailState.Film

	msg := messages.BuildFilmDetailMessage(session, &film) + "\n"

	msg += translator.Translate(session.Lang, "updateChoiceField", nil, nil)

	keyboard := keyboards.BuildFilmUpdateKeyboard(session)

	app.SendImage(film.ImageURL, msg, keyboard)
}

func HandleUpdateFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
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
	if utils.IsCancel(app.Upd) {
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
	msg := translator.Translate(session.Lang, "filmRequestLink", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingURL)
}

func parseUpdateFilmURL(app models.App, session *models.Session) {
	url := utils.ParseMessageString(app.Upd)

	session.FilmDetailState.URL = url

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmImage(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestImage", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingImage)
}

func parseUpdateFilmImage(app models.App, session *models.Session) {
	image, err := utils.ParseImageFromMessage(app.Bot, app.Upd)
	if err != nil {
		msg := translator.Translate(session.Lang, "getImageFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleUpdateFilmCommand(app, session)
		return
	}

	imageURL, err := watchlist.UploadImage(app, image)
	if err != nil {
		msg := translator.Translate(session.Lang, "getImageFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleUpdateFilmCommand(app, session)
		return
	}

	session.FilmDetailState.ImageURL = imageURL

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmTitle(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestTitle", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingTitle)
}

func parseUpdateFilmTitle(app models.App, session *models.Session) {
	session.FilmDetailState.Title = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmDescription(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestDescription", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingDescription)
}

func parseUpdateFilmDescription(app models.App, session *models.Session) {
	session.FilmDetailState.Description = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmGenre(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestGenre", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingGenre)
}

func parseUpdateFilmGenre(app models.App, session *models.Session) {
	session.FilmDetailState.Genre = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmRating(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestRating", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingRating)
}

func parseUpdateFilmRating(app models.App, session *models.Session) {
	session.FilmDetailState.Rating = utils.ParseMessageFloat(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmYear(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestYear", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingYear)
}

func parseUpdateFilmYear(app models.App, session *models.Session) {
	session.FilmDetailState.Year = utils.ParseMessageInt(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmComment(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestComment", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingComment)
}

func parseUpdateFilmComment(app models.App, session *models.Session) {
	session.FilmDetailState.Comment = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmViewed(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestTitle", nil, nil)

	keyboard := keyboards.NewKeyboard().AddSurvey().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingViewed)
}

func parseUpdateFilmViewed(app models.App, session *models.Session) {
	session.FilmDetailState.IsViewed = utils.IsAgree(app.Upd)
	session.FilmDetailState.IsEditViewed = true

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmUserRating(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestUserRating", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingUserRating)
}

func parseUpdateFilmUserRating(app models.App, session *models.Session) {
	session.FilmDetailState.UserRating = utils.ParseMessageFloat(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}

func handleUpdateFilmReview(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filmRequestReview", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingReview)
}

func parseUpdateFilmReview(app models.App, session *models.Session) {
	session.FilmDetailState.Review = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session, HandleUpdateFilmCommand)
}
func updateFilm(app models.App, session *models.Session) {
	film, err := watchlist.UpdateFilm(app, session)
	if err != nil {
		msg := translator.Translate(session.Lang, "updateFilmFailure", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	session.FilmDetailState.Film = *film
}

func finishUpdateFilmProcess(app models.App, session *models.Session, backFunc func(models.App, *models.Session)) {
	state := session.FilmDetailState
	film := state.Film

	if !state.IsEditViewed {
		if state.UserRating == 0 {
			state.UserRating = film.UserRating
		}

		if state.Review == "" {
			state.Review = film.Review
		}

		if !state.IsViewed {
			state.IsViewed = film.IsViewed
		}
	}

	updateFilm(app, session)

	if _, err := GetFilms(app, session); err != nil {
		msg := translator.Translate(session.Lang, "updateFilmListFailure", nil, nil)
		app.SendMessage(msg, nil)
	}

	session.ClearAllStates()
	backFunc(app, session)
}
