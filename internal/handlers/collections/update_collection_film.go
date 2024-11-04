package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
)

var updateCollectionFilmButtons = []builders.Button{
	{"Изображение", states.CallbackUpdateCollectionFilmSelectImage},
	{"Название", states.CallbackUpdateCollectionFilmSelectTitle},
	{"Описание", states.CallbackUpdateCollectionFilmSelectDescription},
	{"Жанр", states.CallbackUpdateCollectionFilmSelectGenre},
	{"Рейтинг", states.CallbackUpdateCollectionFilmSelectRating},
	{"Год выпуска", states.CallbackUpdateCollectionFilmSelectYear},
	{"Комментарий", states.CallbackUpdateCollectionFilmSelectComment},
	{"Просмотрено", states.CallbackUpdateCollectionFilmSelectViewed},
}

var updateCollectionFilmsAfterViewedButtons = []builders.Button{
	{"Оценка пользователя", states.CallbackUpdateCollectionFilmSelectUserRating},
	{"Рецензия", states.CallbackUpdateCollectionFilmSelectReview},
}

func HandleUpdateCollectionFilmCommand(app models.App, session *models.Session) {
	film := session.CollectionFilmState.Object

	buttons := updateCollectionFilmButtons

	if film.IsViewed {
		buttons = append(buttons, updateCollectionFilmsAfterViewedButtons...)
	}

	msg := builders.BuildCollectionFilmDetailMessage(&film)
	msg += "Выберите, какое поле вы хотите изменить?"

	keyboard := builders.NewKeyboard(2).
		AddSeveral(buttons).
		AddBack(states.CallbackUpdateCollectionFilmSelectBack).
		Build()

	app.SendImage(film.ImageURL, msg, keyboard)
}

func HandleUpdateCollectionFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackUpdateCollectionFilmSelectBack:
		HandleCollectionFilmsCommand(app, session)

	case states.CallbackUpdateCollectionFilmSelectImage:
		handleUpdateCollectionFilmImage(app, session)

	case states.CallbackUpdateCollectionFilmSelectTitle:
		handleUpdateCollectionFilmTitle(app, session)

	case states.CallbackUpdateCollectionFilmSelectDescription:
		handleUpdateCollectionFilmDescription(app, session)

	case states.CallbackUpdateCollectionFilmSelectGenre:
		handleUpdateCollectionFilmGenre(app, session)

	case states.CallbackUpdateCollectionFilmSelectRating:
		handleUpdateCollectionFilmRating(app, session)

	case states.CallbackUpdateCollectionFilmSelectYear:
		handleUpdateCollectionFilmYear(app, session)

	case states.CallbackUpdateCollectionFilmSelectComment:
		handleUpdateCollectionFilmComment(app, session)

	case states.CallbackUpdateCollectionFilmSelectViewed:
		handleUpdateCollectionFilmViewed(app, session)

	case states.CallbackUpdateCollectionFilmSelectUserRating:
		handleUpdateCollectionFilmUserRating(app, session)

	case states.CallbackUpdateCollectionFilmSelectReview:
		handleUpdateCollectionFilmReview(app, session)
	}
}

func HandleUpdateCollectionFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		HandleUpdateCollectionFilmCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessUpdateCollectionFilmAwaitingImage:
		parseUpdateCollectionFilmImage(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingTitle:
		parseUpdateCollectionFilmTitle(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingDescription:
		parseUpdateCollectionFilmDescription(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingGenre:
		parseUpdateCollectionFilmGenre(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingRating:
		parseUpdateCollectionFilmRating(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingYear:
		parseUpdateCollectionFilmYear(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingComment:
		parseUpdateCollectionFilmComment(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingViewed:
		parseUpdateCollectionFilmViewed(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingUserRating:
		parseUpdateCollectionFilmUserRating(app, session)

	case states.ProcessUpdateCollectionFilmAwaitingReview:
		parseUpdateCollectionFilmReview(app, session)
	}
}

func handleUpdateCollectionFilmImage(app models.App, session *models.Session) {
	msg := "Введите новую ссылку на изображение фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingImage)
}

func parseUpdateCollectionFilmImage(app models.App, session *models.Session) {
	session.CollectionFilmState.ImageURL = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmTitle(app models.App, session *models.Session) {
	msg := "Введите новое название фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingTitle)
}

func parseUpdateCollectionFilmTitle(app models.App, session *models.Session) {
	session.CollectionFilmState.Title = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmDescription(app models.App, session *models.Session) {
	msg := "Введите новое описание фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingDescription)
}

func parseUpdateCollectionFilmDescription(app models.App, session *models.Session) {
	session.CollectionFilmState.Description = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmGenre(app models.App, session *models.Session) {
	msg := "Введите новый жанр фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingGenre)
}

func parseUpdateCollectionFilmGenre(app models.App, session *models.Session) {
	session.CollectionFilmState.Genre = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmRating(app models.App, session *models.Session) {
	msg := "Введите новый рейтинг фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingRating)
}

func parseUpdateCollectionFilmRating(app models.App, session *models.Session) {
	session.CollectionFilmState.Rating = utils.ParseMessageFloat(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmYear(app models.App, session *models.Session) {
	msg := "Введите новый год выпуска"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingYear)
}

func parseUpdateCollectionFilmYear(app models.App, session *models.Session) {
	session.CollectionFilmState.Year = utils.ParseMessageInt(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmComment(app models.App, session *models.Session) {
	msg := "Введите новый комментарий к фильму"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingComment)
}

func parseUpdateCollectionFilmComment(app models.App, session *models.Session) {
	session.CollectionFilmState.Comment = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmViewed(app models.App, session *models.Session) {
	msg := "Вы просмотрели этот фильм?"

	keyboard := builders.NewKeyboard(2).AddSurvey().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingViewed)
}

func parseUpdateCollectionFilmViewed(app models.App, session *models.Session) {
	session.CollectionFilmState.IsViewed = utils.IsAgree(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmUserRating(app models.App, session *models.Session) {
	msg := "Введите новый пользовательский рейтинг"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingUserRating)
}

func parseUpdateCollectionFilmUserRating(app models.App, session *models.Session) {
	session.CollectionFilmState.UserRating = utils.ParseMessageFloat(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func handleUpdateCollectionFilmReview(app models.App, session *models.Session) {
	msg := "Введите новую рецензию на фильм"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionFilmAwaitingReview)
}

func parseUpdateCollectionFilmReview(app models.App, session *models.Session) {
	session.CollectionFilmState.Review = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionFilmProcess(app, session)
}

func updateCollectionFilm(app models.App, session *models.Session) {
	film, err := watchlist.UpdateFilm(app, session)
	if err != nil {
		app.SendMessage("Не удалось обновить фильм в коллекции", nil)
		return
	}
	log.Println(film.IsViewed)
	session.CollectionFilmState.Object = *film
	app.SendMessage("Фильм в коллекции успешно обновлен", nil)
}

func finishUpdateCollectionFilmProcess(app models.App, session *models.Session) {
	state := session.CollectionFilmState
	collectionFilm := state.Object

	if state.ImageURL == "" {
		state.ImageURL = collectionFilm.ImageURL
	}

	if state.Title == "" {
		state.Title = collectionFilm.Title
	}

	if state.Description == "" {
		state.Description = collectionFilm.Description
	}

	if state.Genre == "" {
		state.Genre = collectionFilm.Genre
	}

	if state.Rating == 0 {
		state.Rating = collectionFilm.Rating
	}

	if state.Year == 0 {
		state.Year = collectionFilm.Year
	}

	if state.Comment == "" {
		state.Comment = collectionFilm.Comment
	}

	if state.UserRating == 0 {
		state.UserRating = collectionFilm.UserRating
	}

	if state.Review == "" {
		state.Review = collectionFilm.Review
	}

	log.Println(state.IsViewed)

	updateCollectionFilm(app, session)
	session.CollectionFilmState.Clear()
	session.ClearState()
	HandleUpdateCollectionFilmCommand(app, session)
}
