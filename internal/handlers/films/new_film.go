package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleNewFilmCommand(app models.App, session *models.Session) {
	keyboard := builders.NewKeyboard(1).AddCancel().Build()
	app.SendMessage("Введите название фильма", keyboard)
	session.SetState(states.ProcessNewFilmAwaitingTitle)
}

func HandleNewFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		session.FilmDetailState.Clear()
		HandleFilmsCommand(app, session)
		return
	}

	switch session.State {
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

	case states.ProcessNewFilmAwaitingViewed:
		parseNewFilmViewed(app, session)

	case states.ProcessNewFilmAwaitingUserRating:
		parseNewFilmUserRating(app, session)

	case states.ProcessNewFilmAwaitingReview:
		parseNewFilmReview(app, session)
	}
}

func parseNewFilmTitle(app models.App, session *models.Session) {
	session.FilmDetailState.Title = utils.ParseMessageString(app.Upd)

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите год выпуска фильма"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingYear)
}

func parseNewFilmYear(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Year = 0
	} else {
		session.FilmDetailState.Year = utils.ParseMessageInt(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите жанр фильма"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingGenre)
}

func parseNewFilmGenre(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Genre = ""
	} else {
		session.FilmDetailState.Genre = utils.ParseMessageString(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите описание"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingDescription)
}

func parseNewFilmDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Description = ""
	} else {
		session.FilmDetailState.Description = utils.ParseMessageString(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите рейтинг фильма"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingRating)
}

func parseNewFilmRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Rating = 0
	} else {
		session.FilmDetailState.Rating = utils.ParseMessageFloat(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Отправьте изображение или ссылку на него"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingImage)
}

func parseNewFilmImage(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		requestNewFilmComment(app, session)
		return
	}

	image, err := utils.ParseServerImage(app.Bot, app.Upd, app.Vars.Host)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		requestNewFilmComment(app, session)
		return
	}

	imageURL, err := watchlist.UploadImage(app, image)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		requestNewFilmComment(app, session)
		return
	}

	session.FilmDetailState.ImageURL = imageURL
	requestNewFilmComment(app, session)
}

func requestNewFilmComment(app models.App, session *models.Session) {
	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите ваш комментарий к фильму"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingComment)
}

func parseNewFilmComment(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Comment = ""
	} else {
		session.FilmDetailState.Comment = utils.ParseMessageString(app.Upd)
	}

	keyboard := builders.NewKeyboard(2).AddSurvey().AddCancel().Build()
	msg := "Вы уже смотрели этот фильм?"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingViewed)
}

func parseNewFilmViewed(app models.App, session *models.Session) {
	switch utils.IsAgree(app.Upd) {
	case true:
		session.FilmDetailState.IsViewed = true
		keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
		msg := "Укажите вашу оценку фильму"
		app.SendMessage(msg, keyboard)
		session.SetState(states.ProcessNewFilmAwaitingUserRating)

	case false:
		session.FilmDetailState.IsViewed = false
		session.FilmDetailState.UserRating = 0
		session.FilmDetailState.Review = ""
		createNewFilm(app, session)
		session.ClearState()
	}
}

func parseNewFilmUserRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.UserRating = 0
	} else {
		session.FilmDetailState.UserRating = utils.ParseMessageFloat(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите ваш отзыв к фильму"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingReview)
}

func parseNewFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.FilmDetailState.Review = ""
	} else {
		session.FilmDetailState.Review = utils.ParseMessageString(app.Upd)
	}

	createNewFilm(app, session)
	session.FilmDetailState.Clear()
	session.ClearState()
}

func createNewFilm(app models.App, session *models.Session) {
	film, err := watchlist.CreateFilm(app, session)
	if err != nil {
		app.SendMessage("Не удалось создать фильм", nil)
		return
	}

	msg := "Новый фильм успешно создан\n"

	msg += builders.BuildFilmDetailMessage(film)

	imageURL := film.ImageURL
	app.SendImage(imageURL, msg, nil)
}
