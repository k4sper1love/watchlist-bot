package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
)

var updateFilmButtons = []builders.Button{
	{"Изображение", states.CallbackUpdateFilmSelectImage},
	{"Название", states.CallbackUpdateFilmSelectTitle},
	{"Описание", states.CallbackUpdateFilmSelectDescription},
	{"Жанр", states.CallbackUpdateFilmSelectGenre},
	{"Рейтинг", states.CallbackUpdateFilmSelectRating},
	{"Год выпуска", states.CallbackUpdateFilmSelectYear},
	{"Комментарий", states.CallbackUpdateFilmSelectComment},
	{"Просмотрено", states.CallbackUpdateFilmSelectViewed},
}

var updateFilmsAfterViewedButtons = []builders.Button{
	{"Оценка пользователя", states.CallbackUpdateFilmSelectUserRating},
	{"Рецензия", states.CallbackUpdateFilmSelectReview},
}

func HandleUpdateFilmCommand(app models.App, session *models.Session) {
	film := session.FilmDetailState.Object

	buttons := updateFilmButtons

	if film.IsViewed {
		buttons = append(buttons, updateFilmsAfterViewedButtons...)
	}

	msg := builders.BuildFilmDetailMessage(&film)
	msg += "Выберите, какое поле вы хотите изменить?"

	keyboard := builders.NewKeyboard(2).
		AddSeveral(buttons).
		AddBack(states.CallbackUpdateFilmSelectBack).
		Build()

	app.SendImage(film.ImageURL, msg, keyboard)
}

func HandleUpdateFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackUpdateFilmSelectBack:
		HandleManageFilmCommand(app, session)

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
		session.ClearState()
		session.FilmDetailState.Clear()
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

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingImage)
}

func parseUpdateFilmImage(app models.App, session *models.Session) {
	image, err := utils.ParseServerImage(app.Bot, app.Upd, app.Vars.Host)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		finishUpdateFilmProcess(app, session)
		HandleUpdateFilmCommand(app, session)
		return
	}

	imageURL, err := watchlist.UploadImage(app, image)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		finishUpdateFilmProcess(app, session)
		HandleUpdateFilmCommand(app, session)
		return
	}

	session.FilmDetailState.ImageURL = imageURL

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmTitle(app models.App, session *models.Session) {
	msg := "Введите новое название фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingTitle)
}

func parseUpdateFilmTitle(app models.App, session *models.Session) {
	session.FilmDetailState.Title = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmDescription(app models.App, session *models.Session) {
	msg := "Введите новое описание фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingDescription)
}

func parseUpdateFilmDescription(app models.App, session *models.Session) {
	session.FilmDetailState.Description = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmGenre(app models.App, session *models.Session) {
	msg := "Введите новый жанр фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingGenre)
}

func parseUpdateFilmGenre(app models.App, session *models.Session) {
	session.FilmDetailState.Genre = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmRating(app models.App, session *models.Session) {
	msg := "Введите новый рейтинг фильма"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingRating)
}

func parseUpdateFilmRating(app models.App, session *models.Session) {
	session.FilmDetailState.Rating = utils.ParseMessageFloat(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmYear(app models.App, session *models.Session) {
	msg := "Введите новый год выпуска"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingYear)
}

func parseUpdateFilmYear(app models.App, session *models.Session) {
	session.FilmDetailState.Year = utils.ParseMessageInt(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmComment(app models.App, session *models.Session) {
	msg := "Введите новый комментарий к фильму"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingComment)
}

func parseUpdateFilmComment(app models.App, session *models.Session) {
	session.FilmDetailState.Comment = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmViewed(app models.App, session *models.Session) {
	msg := "Вы посмотрели этот фильм?"

	keyboard := builders.NewKeyboard(2).AddSurvey().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingViewed)
}

func parseUpdateFilmViewed(app models.App, session *models.Session) {
	session.FilmDetailState.IsViewed = utils.IsAgree(app.Upd)
	session.FilmDetailState.IsEditViewed = true

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmUserRating(app models.App, session *models.Session) {
	msg := "Введите новый пользовательский рейтинг"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingUserRating)
}

func parseUpdateFilmUserRating(app models.App, session *models.Session) {
	session.FilmDetailState.UserRating = utils.ParseMessageFloat(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func handleUpdateFilmReview(app models.App, session *models.Session) {
	msg := "Введите новую рецензию на фильм"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateFilmAwaitingReview)
}

func parseUpdateFilmReview(app models.App, session *models.Session) {
	session.FilmDetailState.Review = utils.ParseMessageString(app.Upd)

	finishUpdateFilmProcess(app, session)
	HandleUpdateFilmCommand(app, session)
}

func updateFilm(app models.App, session *models.Session) {
	film, err := watchlist.UpdateFilm(app, session)
	if err != nil {
		log.Println(err)
		app.SendMessage("Не удалось обновить фильм", nil)
		return
	}

	session.FilmDetailState.Object = *film
}

func finishUpdateFilmProcess(app models.App, session *models.Session) {
	state := session.FilmDetailState
	film := state.Object

	if state.ImageURL == "" {
		state.ImageURL = film.ImageURL
	}

	if state.Title == "" {
		state.Title = film.Title
	}

	if state.Description == "" {
		state.Description = film.Description
	}

	if state.Genre == "" {
		state.Genre = film.Genre
	}

	if state.Rating == 0 {
		state.Rating = film.Rating
	}

	if state.Year == 0 {
		state.Year = film.Year
	}

	if state.Comment == "" {
		state.Comment = film.Comment
	}

	if !state.IsViewed && !state.IsEditViewed {
		state.IsViewed = film.IsViewed
	}

	if state.UserRating == 0 {
		state.UserRating = film.UserRating
	}

	if state.Review == "" {
		state.Review = film.Review
	}

	updateFilm(app, session)

	if _, err := getFilms(app, session); err != nil {
		app.SendMessage("Ошибка при обновлении списка фильмов", nil)
	}

	session.FilmDetailState.Clear()
	session.ClearState()
}
