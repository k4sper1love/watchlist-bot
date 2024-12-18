package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/parsing"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleNewFilmCommand(app models.App, session *models.Session) {
	keyboard := keyboards.BuildFilmNewKeyboard()
	app.SendMessage("Выберите один из предложенных методов", keyboard)
}

func HandleNewFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleNewFilmCommand(app, session)
		return
	}

	switch session.State {
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

	case states.ProcessNewFilmAwaitingViewed:
		parseNewFilmViewed(app, session)

	case states.ProcessNewFilmAwaitingUserRating:
		parseNewFilmUserRating(app, session)

	case states.ProcessNewFilmAwaitingReview:
		parseNewFilmReview(app, session)
	}
}

func HandleNewFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackNewFilmSelectBack:
		HandleFilmsCommand(app, session)

	case states.CallbackNewFilmSelectManually:
		handleNewFilmManually(app, session)

	case states.CallbackNewFilmSelectFromURL:
		handleNewFilmFromURL(app, session)
	}
}

func handleNewFilmFromURL(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddCancel().Build()
	app.SendMessage("Пришлите ссылку на фильм (kinopoisk, kinopoisk HD, rezka)", keyboard)
	session.SetState(states.ProcessNewFilmAwaitingURL)
}

func parseNewFilmFromURL(app models.App, session *models.Session) {
	url := utils.ParseMessageString(app.Upd)

	film, err := parsing.GetFilmByURL(app, url)
	if err != nil {
		app.SendMessage("Ошибка при получении фильма", nil)
		session.ClearAllStates()
		HandleNewFilmCommand(app, session)
		return
	}

	session.FilmDetailState.SetFromFilm(film)

	imageURL, err := parseAndUploadImageFromURL(app, film.ImageURL)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		session.FilmDetailState.SetImageURL("")
		requestNewFilmComment(app, session)
		return
	}
	session.FilmDetailState.SetImageURL(imageURL)

	requestNewFilmComment(app, session)
}

func handleNewFilmManually(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddCancel().Build()
	app.SendMessage("Введите название фильма", keyboard)
	session.SetState(states.ProcessNewFilmAwaitingTitle)
}

func parseNewFilmTitle(app models.App, session *models.Session) {
	session.FilmDetailState.Title = utils.ParseMessageString(app.Upd)

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
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

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
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

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
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

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
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

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
	msg := "Отправьте изображение или ссылку на него"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingImage)
}

func parseNewFilmImage(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		requestNewFilmComment(app, session)
		return
	}

	imageURL, err := parseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage("Ошибка при получении изображения", nil)
		requestNewFilmComment(app, session)
		return
	}

	session.FilmDetailState.SetImageURL(imageURL)

	requestNewFilmComment(app, session)
}

func requestNewFilmComment(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
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

	keyboard := keyboards.NewKeyboard().AddSurvey().AddCancel().Build()
	msg := "Вы уже смотрели этот фильм?"

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewFilmAwaitingViewed)
}

func parseNewFilmViewed(app models.App, session *models.Session) {
	switch utils.IsAgree(app.Upd) {
	case true:
		session.FilmDetailState.IsViewed = true
		keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
		msg := "Укажите вашу оценку фильму"
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

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build()
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

	finishNewFilmProcess(app, session)
}

func finishNewFilmProcess(app models.App, session *models.Session) {
	if err := CreateNewFilm(app, session); err != nil {
		app.SendMessage("Не удалось создать фильм", nil)
		session.ClearAllStates()
		HandleFilmsCommand(app, session)
		return
	}
	app.SendMessage("Новый фильм успешно создан", nil)

	session.ClearAllStates()

	totalRecords := session.FilmsState.TotalRecords
	pageSize := session.FilmsState.PageSize

	newPage, newIndex := utils.CalculateNewElementPageAndIndex(totalRecords, pageSize)

	session.FilmsState.CurrentPage = newPage
	session.FilmDetailState.Index = newIndex

	if _, err := GetFilms(app, session); err != nil {
		app.SendMessage("Ошибка при обновлении списка фильмов", nil)
	}

	HandleFilmsDetailCommand(app, session)
}

func CreateNewFilm(app models.App, session *models.Session) error {
	switch session.Context {
	case states.ContextFilm:
		return createNewUserFilm(app, session)

	case states.ContextCollection:
		return createNewCollectionFilm(app, session)

	default:
		return fmt.Errorf("unsupported session context: %v", session.Context)
	}
}

func createNewUserFilm(app models.App, session *models.Session) error {
	_, err := watchlist.CreateFilm(app, session)
	if err != nil {
		return err
	}

	return nil
}

func createNewCollectionFilm(app models.App, session *models.Session) error {
	collectionFilm, err := watchlist.CreateCollectionFilm(app, session)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Новый фильм в коллекции %q успешно создан\n", collectionFilm.Collection.Name)

	msg += messages.BuildFilmDetailMessage(&collectionFilm.Film)

	imageURL := collectionFilm.Film.ImageURL
	app.SendImage(imageURL, msg, nil)

	return nil
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
