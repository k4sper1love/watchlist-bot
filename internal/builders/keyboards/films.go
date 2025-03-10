package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateFilmButtons = []Button{
	{"", "image", states.CallbackUpdateFilmSelectImage, "", true},
	{"", "title", states.CallbackUpdateFilmSelectTitle, "", true},
	{"", "description", states.CallbackUpdateFilmSelectDescription, "", true},
	{"", "genre", states.CallbackUpdateFilmSelectGenre, "", true},
	{"", "rating", states.CallbackUpdateFilmSelectRating, "", true},
	{"", "yearOfRelease", states.CallbackUpdateFilmSelectYear, "", true},
	{"", "comment", states.CallbackUpdateFilmSelectComment, "", true},
	{"", "viewed", states.CallbackUpdateFilmSelectViewed, "", true},
}

var updateFilmAfterViewedButtons = []Button{
	{"", "userRating", states.CallbackUpdateFilmSelectUserRating, "", true},
	{"", "review", states.CallbackUpdateFilmSelectReview, "", true},
}

func Films(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(len(session.FilmsState.Films) > 0, func(k *Keyboard) {
			k.AddSearch(states.CallbackFilmsFind)
		}).
		AddFilmSelect(session).
		AddNavigation(currentPage, lastPage, states.PrefixFilmsPage, true).
		AddFilmFiltersAndSorting(session).
		AddIf(session.Context == states.ContextFilm, func(k *Keyboard) {
			k.AddFilmNew()
			k.AddBack("")
		}).
		AddIf(session.Context == states.ContextCollection, func(k *Keyboard) {
			k.AddCollectionFilmFromCollection()
			k.AddFavorite(session.CollectionDetailState.Collection.IsFavorite, states.CallbackCollectionsFavorite)
			k.AddManage(states.CallbackCollectionsManage)
			k.AddBack(states.CallbackFilmsBack)
		}).
		Build(session.Lang)
}

func FindFilms(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFilmSelect(session).
		AddNavigation(currentPage, lastPage, states.PrefixFindFilmsPage, true).
		AddBack(states.CallbackFindFilmsBack).
		Build(session.Lang)
}

func FindNewFilm(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFindNewFilmSelect(session).
		AddNavigation(currentPage, lastPage, states.PrefixFindNewFilmPage, true).
		AddBack(states.CallbackFindNewFilmBack).
		Build(session.Lang)
}

func FilmDetail(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	film := session.FilmDetailState.Film
	return New().
		AddFavorite(film.IsFavorite, states.CallbackFilmDetailFavorite).
		AddIf(film.URL != "", func(k *Keyboard) {
			k.AddOpenInBrowser(film.URL)
		}).
		AddIf(!film.IsViewed, func(k *Keyboard) {
			k.AddFilmViewed()
		}).
		AddManage(states.CallbackFilmsManage).
		AddIf(session.Context == states.ContextFilm, func(k *Keyboard) {
			k.AddCollectionFilmFromFilm()
		}).
		AddIf(session.FilmDetailState.HasIndex(), func(k *Keyboard) {
			itemID := utils.GetItemID(session.FilmDetailState.Index, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
			k.AddNavigation(itemID, session.FilmsState.TotalRecords, states.PrefixFilmDetailPage, false)
		}).
		AddBack(states.CallbackFilmDetailBack).
		Build(session.Lang)
}

func FilmManage(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.Context == states.ContextCollection, func(k *Keyboard) {
			k.AddFilmRemoveFromCollection()
		}).
		AddUpdate(states.CallbackManageFilmSelectUpdate).
		AddDelete(states.CallbackManageFilmSelectDelete).
		AddBack(states.CallbackManageFilmSelectBack).
		Build(session.Lang)
}

func FilmNew(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddNewFilmManually().
		AddNewFilmFromURL().
		AddNewFilmFind().
		AddBack(states.CallbackNewFilmSelectBack).
		Build(session.Lang)
}

func FilmUpdate(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddUpdateFilmURL().
		AddButtonsWithRowSize(2, updateFilmButtons...).
		AddIf(session.FilmDetailState.Film.IsViewed, func(k *Keyboard) {
			k.AddButtonsWithRowSize(2, updateFilmAfterViewedButtons...)
		}).
		AddBack(states.CallbackUpdateFilmSelectBack).
		Build(session.Lang)
}

func FilmsFilter(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	filter := session.GetFilmsFiltersByContext()
	return New().
		AddButtons(getFiltersFilmsButtons(filter, session.Lang)...).
		AddIf(filter.IsFiltersEnabled(), func(k *Keyboard) {
			k.AddResetAllFilmsFilters()
		}).
		AddBack(states.CallbackFiltersFilmsBack).
		Build(session.Lang)
}

func FilmsSorting(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	sorting := session.GetFilmsSortingByContext()
	return New().
		AddButtons(getSortingFilmsButtons(sorting, session.Lang)...).
		AddIf(sorting.IsSortingEnabled(), func(k *Keyboard) {
			k.AddResetAllSorting(states.CallbackSortingFilmsAllReset)
		}).
		AddBack(states.CallbackSortingFilmsBack).
		Build(session.Lang)
}

func FilmsNotFound(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddAgain(states.CallbackFindFilmsAgain).
		AddBack(states.CallbackFindFilmsBack).
		Build(session.Lang)
}

func FindNewFilmsNotFound(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddAgain(states.CallbackFindNewFilmAgain).
		AddBack(states.CallbackFindNewFilmBack).
		Build(session.Lang)
}

func NewFilmChangeToken(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddChangeToken().
		AddBack(states.CallbackFilmsNew).
		Build(session.Lang)
}

func FilmFilterSwitch(session *models.Session, filterType string) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSurvey().
		AddResetFilmsFilter(session, filterType).
		AddCancel().
		Build(session.Lang)
}

func FilmFilterRange(session *models.Session, filterType string) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddResetFilmsFilter(session, filterType).
		AddCancel().
		Build(session.Lang)
}
