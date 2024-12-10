package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateFilmButtons = []keyboards.Button{
	{"Изображение", states.CallbackUpdateFilmSelectImage},
	{"Название", states.CallbackUpdateFilmSelectTitle},
	{"Описание", states.CallbackUpdateFilmSelectDescription},
	{"Жанр", states.CallbackUpdateFilmSelectGenre},
	{"Рейтинг", states.CallbackUpdateFilmSelectRating},
	{"Год выпуска", states.CallbackUpdateFilmSelectYear},
	{"Комментарий", states.CallbackUpdateFilmSelectComment},
	{"Просмотрено", states.CallbackUpdateFilmSelectViewed},
}

var updateFilmsAfterViewedButtons = []keyboards.Button{
	{"Оценка пользователя", states.CallbackUpdateFilmSelectUserRating},
	{"Рецензия", states.CallbackUpdateFilmSelectReview},
}

func HandleUpdateFilmCommand(app models.App, session *models.Session) {
	film := session.FilmDetailState.Film

	msg := messages.BuildFilmDetailMessage(&film)
	msg += "\nВыберите, какое поле вы хотите изменить?"

	keyboard := keyboards.BuildFilmUpdateKeyboard(session)

	app.SendImage(film.ImageURL, msg, keyboard)
}

func HandleUpdateFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackUpdateFilmSelectBack:
		HandleFilmsDetailCommand(app, session)

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

func handleUpdateFilmImage(app models.App, session *models.Session) {
	msg := "Отправьте новое изображение или ссылку на него"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingImage)
}

func parseUpdateFilmImage(app models.App, session *models.Session) {
	image, err := utils.ParseImageFromMessage(app.Bot, app.Upd)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		session.ClearAllStates()
		HandleUpdateFilmCommand(app, session)
		return
	}

	imageURL, err := watchlist.UploadImage(app, image)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		session.ClearAllStates()
		HandleUpdateFilmCommand(app, session)
		return
	}

	session.FilmDetailState.ImageURL = imageURL

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmTitle(app models.App, session *models.Session) {
	msg := "Введите новое название фильма"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingTitle)
}

func parseUpdateFilmTitle(app models.App, session *models.Session) {
	session.FilmDetailState.Title = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmDescription(app models.App, session *models.Session) {
	msg := "Введите новое описание фильма"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingDescription)
}

func parseUpdateFilmDescription(app models.App, session *models.Session) {
	session.FilmDetailState.Description = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmGenre(app models.App, session *models.Session) {
	msg := "Введите новый жанр фильма"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingGenre)
}

func parseUpdateFilmGenre(app models.App, session *models.Session) {
	session.FilmDetailState.Genre = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmRating(app models.App, session *models.Session) {
	msg := "Введите новый рейтинг фильма"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingRating)
}

func parseUpdateFilmRating(app models.App, session *models.Session) {
	session.FilmDetailState.Rating = utils.ParseMessageFloat(app.Upd)

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmYear(app models.App, session *models.Session) {
	msg := "Введите новый год выпуска"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingYear)
}

func parseUpdateFilmYear(app models.App, session *models.Session) {
	session.FilmDetailState.Year = utils.ParseMessageInt(app.Upd)

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmComment(app models.App, session *models.Session) {
	msg := "Введите новый комментарий к фильму"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingComment)
}

func parseUpdateFilmComment(app models.App, session *models.Session) {
	session.FilmDetailState.Comment = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmViewed(app models.App, session *models.Session) {
	msg := "Вы посмотрели этот фильм?"

	keyboard := keyboards.NewKeyboard().AddSurvey().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingViewed)
}

func parseUpdateFilmViewed(app models.App, session *models.Session) {
	session.FilmDetailState.IsViewed = utils.IsAgree(app.Upd)
	session.FilmDetailState.IsEditViewed = true

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmUserRating(app models.App, session *models.Session) {
	msg := "Введите новый пользовательский рейтинг"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingUserRating)
}

func parseUpdateFilmUserRating(app models.App, session *models.Session) {
	session.FilmDetailState.UserRating = utils.ParseMessageFloat(app.Upd)

	finishUpdateFilmProcess(app, session)
}

func handleUpdateFilmReview(app models.App, session *models.Session) {
	msg := "Введите новую рецензию на фильм"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingReview)
}

func parseUpdateFilmReview(app models.App, session *models.Session) {
	session.FilmDetailState.Review = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
}
func updateFilm(app models.App, session *models.Session) {
	film, err := watchlist.UpdateFilm(app, session)
	if err != nil {
		app.SendMessage("Не удалось обновить фильм", nil)
		return
	}

	session.FilmDetailState.Film = *film
}

func finishUpdateFilmProcess(app models.App, session *models.Session) {
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
		app.SendMessage("Ошибка при обновлении списка фильмов", nil)
	}

	session.ClearAllStates()
	HandleUpdateFilmCommand(app, session)
}
