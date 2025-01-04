package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateFilmURLButton = Button{"", "filmURL", states.CallbackUpdateFilmSelectURL, ""}

var updateFilmButtons = []Button{
	{"", "image", states.CallbackUpdateFilmSelectImage, ""},
	{"", "title", states.CallbackUpdateFilmSelectTitle, ""},
	{"", "description", states.CallbackUpdateFilmSelectDescription, ""},
	{"", "genre", states.CallbackUpdateFilmSelectGenre, ""},
	{"", "rating", states.CallbackUpdateFilmSelectRating, ""},
	{"", "yearOfRelease", states.CallbackUpdateFilmSelectYear, ""},
	{"", "comment", states.CallbackUpdateFilmSelectComment, ""},
	{"", "viewed", states.CallbackUpdateFilmSelectViewed, ""},
}

var updateFilmsAfterViewedButtons = []Button{
	{"", "userRating", states.CallbackUpdateFilmSelectUserRating, ""},
	{"", "Review", states.CallbackUpdateFilmSelectReview, ""},
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

	return keyboard.Build(session.Lang)
}

func BuildFilmDetailKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	index := session.FilmDetailState.Index
	itemID := utils.GetItemID(index, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
	film := session.FilmDetailState.Film

	keyboard := NewKeyboard()

	if film.URL != "" {
		keyboard.AddURLButton("", "openInBrowser", film.URL)
	}

	if !film.IsViewed {
		keyboard.AddButton("✔️", "viewed", states.CallbackFilmDetailViewed, "")
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

	return keyboard.Build(session.Lang)
}

func BuildFilmManageKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFilmUpdate()

	if session.Context == states.ContextCollection {
		keyboard.AddFilmRemoveFromCollection()
	}

	keyboard.AddFilmDelete()

	keyboard.AddBack(states.CallbackManageFilmSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmNewKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddNewFilmManually()

	keyboard.AddNewFilmFromURL()

	keyboard.AddBack(states.CallbackNewFilmSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmUpdateKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	film := session.FilmDetailState.Film
	keyboard := NewKeyboard()

	keyboard.AddRow(updateFilmURLButton)

	keyboard.AddButtonsWithRowSize(2, updateFilmButtons...)

	if film.IsViewed {
		keyboard.AddButtonsWithRowSize(2, updateFilmsAfterViewedButtons...)
	}

	keyboard.AddBack(states.CallbackUpdateFilmSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmViewedKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddSkip()

	keyboard.AddCancel()

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)

		buttons = append(buttons, Button{"", fmt.Sprintf("%s (%d)", film.Title, itemID), fmt.Sprintf("select_film_%d", i), ""})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddFilmNew() *Keyboard {
	return k.AddButton("➕", "createFilm", states.CallbackFilmsNew, "")
}

func (k *Keyboard) AddFilmDelete() *Keyboard {
	return k.AddButton("🗑️", "deleteFilm", states.CallbackManageFilmSelectDelete, "")
}

func (k *Keyboard) AddFilmUpdate() *Keyboard {
	return k.AddButton("✏️", "updateFilm", states.CallbackManageFilmSelectUpdate, "")
}

func (k *Keyboard) AddFilmManage() *Keyboard {
	return k.AddButton("⚙️", "manageFilm", states.CallbackFilmsManage, "")
}

func (k *Keyboard) AddFilmRemoveFromCollection() *Keyboard {
	return k.AddButton("❌", "removeFilmFromCollection", states.CallbackManageFilmSelectRemoveFromCollection, "")
}

func (k *Keyboard) AddNewFilmManually() *Keyboard {
	return k.AddButton("", "manually", states.CallbackNewFilmSelectManually, "")
}

func (k *Keyboard) AddNewFilmFromURL() *Keyboard {
	return k.AddButton("", "fromURL", states.CallbackNewFilmSelectFromURL, "")
}
