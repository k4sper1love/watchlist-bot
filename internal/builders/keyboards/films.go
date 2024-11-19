package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateFilmButtons = []Button{
	{"Изображение", states.CallbackUpdateFilmSelectImage},
	{"Название", states.CallbackUpdateFilmSelectTitle},
	{"Описание", states.CallbackUpdateFilmSelectDescription},
	{"Жанр", states.CallbackUpdateFilmSelectGenre},
	{"Рейтинг", states.CallbackUpdateFilmSelectRating},
	{"Год выпуска", states.CallbackUpdateFilmSelectYear},
	{"Комментарий", states.CallbackUpdateFilmSelectComment},
	{"Просмотрено", states.CallbackUpdateFilmSelectViewed},
}

var updateFilmsAfterViewedButtons = []Button{
	{"Оценка пользователя", states.CallbackUpdateFilmSelectUserRating},
	{"Рецензия", states.CallbackUpdateFilmSelectReview},
}

func BuildFilmsKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFilmSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackFilmsPrevPage,
		states.CallbackFilmsNextPage,
	)

	switch session.Context {
	case states.ContextFilm:
		keyboard.AddFilmNew()
		keyboard.AddBack("")

	case states.ContextCollection:
		keyboard.AddCollectionFilmFromCollection()
		keyboard.AddCollectionsManage()
		keyboard.AddBack(states.CallbackFilmsBack)
	}

	return keyboard.Build()
}

func BuildFilmDetailKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	index := session.FilmDetailState.Index
	itemID := utils.GetItemID(index, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
	film := session.FilmDetailState.Film

	keyboard := NewKeyboard()

	if !film.IsViewed {
		keyboard.AddButton("Просмотрено✔️", states.CallbackFilmDetailViewed)
	}

	keyboard.AddFilmManage()

	if session.Context == states.ContextFilm {
		keyboard.AddCollectionFilmFromFilm()
	}

	keyboard.AddNavigation(
		itemID,
		session.FilmsState.TotalRecords,
		states.CallbackFilmDetailPrevPage,
		states.CallbackFilmDetailNextPage)

	keyboard.AddBack(states.CallbackFilmDetailBack)

	return keyboard.Build()
}

func BuildFilmManageKeyboard() *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFilmUpdate()

	keyboard.AddFilmDelete()

	keyboard.AddBack(states.CallbackManageFilmSelectBack)

	return keyboard.Build()
}

func BuildFilmUpdateKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	film := session.FilmDetailState.Film
	keyboard := NewKeyboard()

	keyboard.AddButtonsWithRowSize(2, updateFilmButtons...)

	if film.IsViewed {
		keyboard.AddButtonsWithRowSize(2, updateFilmsAfterViewedButtons...)
	}

	keyboard.AddBack(states.CallbackUpdateFilmSelectBack)

	return keyboard.Build()
}

func BuildFilmViewedKeyboard() *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddSkip()

	keyboard.AddCancel()

	return keyboard.Build()
}

func (k *Keyboard) AddFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)

		buttons = append(buttons, Button{fmt.Sprintf("%s (%d)", film.Title, itemID), fmt.Sprintf("select_film_%d", i)})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddFilmNew() *Keyboard {
	return k.AddButton("➕ Создать фильм", states.CallbackFilmsNew)
}

func (k *Keyboard) AddFilmDelete() *Keyboard {
	return k.AddButton("🗑️ Удалить фильм", states.CallbackManageFilmSelectDelete)
}

func (k *Keyboard) AddFilmUpdate() *Keyboard {
	return k.AddButton("✏️ Обновить фильм", states.CallbackManageFilmSelectUpdate)
}

func (k *Keyboard) AddFilmManage() *Keyboard {
	return k.AddButton("⚙️ Управление фильмом", states.CallbackFilmsManage)
}
