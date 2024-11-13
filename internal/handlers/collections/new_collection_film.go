package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
)

func HandleNewCollectionFilmCommand(app models.App, session *models.Session) {
	keyboard := builders.NewKeyboard(1).AddCancel().Build()
	app.SendMessage("Введите название фильма", keyboard)
	session.SetState(states.ProcessNewCollectionFilmAwaitingTitle)
}

func HandleNewCollectionFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		session.CollectionFilmState.Clear()
		HandleCollectionFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessNewCollectionFilmAwaitingTitle:
		parseNewCollectionFilmTitle(app, session)

	case states.ProcessNewCollectionFilmAwaitingYear:
		parseNewCollectionFilmYear(app, session)

	case states.ProcessNewCollectionFilmAwaitingGenre:
		parseNewCollectionFilmGenre(app, session)

	case states.ProcessNewCollectionFilmAwaitingDescription:
		parseNewCollectionFilmDescription(app, session)

	case states.ProcessNewCollectionFilmAwaitingRating:
		parseNewCollectionFilmRating(app, session)

	case states.ProcessNewCollectionFilmAwaitingImage:
		parseNewCollectionFilmImage(app, session)

	case states.ProcessNewCollectionFilmAwaitingComment:
		parseNewCollectionFilmComment(app, session)

	case states.ProcessNewCollectionFilmAwaitingViewed:
		parseNewCollectionFilmViewed(app, session)

	case states.ProcessNewCollectionFilmAwaitingUserRating:
		parseNewCollectionFilmUserRating(app, session)

	case states.ProcessNewCollectionFilmAwaitingReview:
		parseNewCollectionFilmReview(app, session)
	}
}

func parseNewCollectionFilmTitle(app models.App, session *models.Session) {
	session.CollectionFilmState.Title = utils.ParseMessageString(app.Upd)

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите год выпуска фильма"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingYear)
}

func parseNewCollectionFilmYear(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.Year = 0
	} else {
		session.CollectionFilmState.Year = utils.ParseMessageInt(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите жанр фильма"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingGenre)
}

func parseNewCollectionFilmGenre(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.Genre = ""
	} else {
		session.CollectionFilmState.Genre = utils.ParseMessageString(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите описание"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingDescription)
}

func parseNewCollectionFilmDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.Description = ""
	} else {
		session.CollectionFilmState.Description = utils.ParseMessageString(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите рейтинг фильма"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingRating)
}

func parseNewCollectionFilmRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.Rating = 0
	} else {
		session.CollectionFilmState.Rating = utils.ParseMessageFloat(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Отправьте изображение или ссылку на него"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingImage)
}

func parseNewCollectionFilmImage(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		requestNewCollectionFilmComment(app, session)
		return
	}

	image, err := utils.ParseServerImage(app.Bot, app.Upd, app.Vars.Host)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		requestNewCollectionFilmComment(app, session)
		return
	}

	imageURL, err := watchlist.UploadImage(app, image)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		requestNewCollectionFilmComment(app, session)
		return
	}

	session.CollectionFilmState.ImageURL = imageURL
	requestNewCollectionFilmComment(app, session)
}

func requestNewCollectionFilmComment(app models.App, session *models.Session) {
	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите ваш комментарий к фильму"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingComment)
}

func parseNewCollectionFilmComment(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.Comment = ""
	} else {
		session.CollectionFilmState.Comment = utils.ParseMessageString(app.Upd)
	}

	keyboard := builders.NewKeyboard(2).AddSurvey().AddCancel().Build()
	msg := "Вы уже смотрели этот фильм?"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingViewed)
}

func parseNewCollectionFilmViewed(app models.App, session *models.Session) {
	switch utils.IsAgree(app.Upd) {
	case true:
		session.CollectionFilmState.IsViewed = true
		keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
		msg := "Укажите вашу оценку фильму"
		app.SendMessage(msg, keyboard)
		session.SetState(states.ProcessNewCollectionFilmAwaitingUserRating)

	case false:
		session.CollectionFilmState.IsViewed = false
		session.CollectionFilmState.UserRating = 0
		session.CollectionFilmState.Review = ""
		createNewCollectionFilm(app, session)
		session.ClearState()
	}
}

func parseNewCollectionFilmUserRating(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.UserRating = 0
	} else {
		session.CollectionFilmState.UserRating = utils.ParseMessageFloat(app.Upd)
	}

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()
	msg := "Укажите ваш отзыв к фильму"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionFilmAwaitingReview)
}

func parseNewCollectionFilmReview(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionFilmState.Review = ""
	} else {
		session.CollectionFilmState.Review = utils.ParseMessageString(app.Upd)
	}

	createNewCollectionFilm(app, session)
	session.CollectionFilmState.Clear()
	session.ClearState()
}

func createNewCollectionFilm(app models.App, session *models.Session) {
	log.Println(session.CollectionFilmState)
	collectionFilm, err := watchlist.CreateCollectionFilm(app, session)
	if err != nil {
		app.SendMessage("Не удалось создать фильм в коллекцию", nil)
		return
	}

	msg := fmt.Sprintf("Новый фильм в коллекции %q успешно создан\n", collectionFilm.Collection.Name)

	msg += builders.BuildCollectionFilmDetailMessage(&collectionFilm.Film)

	imageURL := collectionFilm.Film.ImageURL
	app.SendImage(imageURL, msg, nil)
}
