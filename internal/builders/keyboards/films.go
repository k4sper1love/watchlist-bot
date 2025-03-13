package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateFilmButtons = []Button{
	{"", "image", states.CallUpdateFilmImage, "", true},
	{"", "title", states.CallUpdateFilmTitle, "", true},
	{"", "description", states.CallUpdateFilmDescription, "", true},
	{"", "genre", states.CallUpdateFilmGenre, "", true},
	{"", "rating", states.CallUpdateFilmRating, "", true},
	{"", "yearOfRelease", states.CallUpdateFilmYear, "", true},
	{"", "comment", states.CallUpdateFilmComment, "", true},
	{"", "viewed", states.CallUpdateFilmViewed, "", true},
}

var updateFilmAfterViewedButtons = []Button{
	{"", "userRating", states.CallUpdateFilmUserRating, "", true},
	{"", "review", states.CallUpdateFilmReview, "", true},
}

func Films(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(len(session.FilmsState.Films) > 0, func(k *Keyboard) {
			k.AddSearch(states.CallFilmsFind)
		}).
		AddFilmSelect(session).
		AddNavigation(currentPage, lastPage, states.FilmsPage, true).
		AddFilmFiltersAndSorting(session).
		AddIf(session.Context == states.CtxFilm, func(k *Keyboard) {
			k.AddFilmNew()
			k.AddBack("")
		}).
		AddIf(session.Context == states.CtxCollection, func(k *Keyboard) {
			k.AddCollectionFilmFromCollection()
			k.AddFavorite(session.CollectionDetailState.Collection.IsFavorite, states.CallCollectionsFavorite)
			k.AddManage(states.CallCollectionsManage)
			k.AddBack(states.CallFilmsBack)
		}).
		Build(session.Lang)
}

func FindFilms(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFilmSelect(session).
		AddNavigation(currentPage, lastPage, states.FindFilmsPage, true).
		AddBack(states.CallFindFilmsBack).
		Build(session.Lang)
}

func FindNewFilm(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFindNewFilmSelect(session).
		AddNavigation(currentPage, lastPage, states.FindNewFilmPage, true).
		AddBack(states.CallFindNewFilmBack).
		Build(session.Lang)
}

func FilmDetail(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	film := session.FilmDetailState.Film
	return New().
		AddFavorite(film.IsFavorite, states.CallFilmDetailFavorite).
		AddIf(film.URL != "", func(k *Keyboard) {
			k.AddOpenInBrowser(film.URL)
		}).
		AddIf(!film.IsViewed, func(k *Keyboard) {
			k.AddFilmViewed()
		}).
		AddManage(states.CallFilmsManage).
		AddIf(session.Context == states.CtxFilm, func(k *Keyboard) {
			k.AddCollectionFilmFromFilm()
		}).
		AddIf(session.FilmDetailState.HasIndex(), func(k *Keyboard) {
			itemID := utils.GetItemID(session.FilmDetailState.Index, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
			k.AddNavigation(itemID, session.FilmsState.TotalRecords, states.FilmDetailPage, false)
		}).
		AddBack(states.CallFilmDetailBack).
		Build(session.Lang)
}

func FilmManage(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.Context == states.CtxCollection, func(k *Keyboard) {
			k.AddFilmRemoveFromCollection()
		}).
		AddUpdate(states.CallManageFilmUpdate).
		AddDelete(states.CallManageFilmDelete).
		AddBack(states.CallManageFilmBack).
		Build(session.Lang)
}

func FilmNew(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddNewFilmManually().
		AddNewFilmFromURL().
		AddNewFilmFind().
		AddBack(states.CallNewFilmBack).
		Build(session.Lang)
}

func FilmUpdate(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddUpdateFilmURL().
		AddButtonsWithRowSize(2, updateFilmButtons...).
		AddIf(session.FilmDetailState.Film.IsViewed, func(k *Keyboard) {
			k.AddButtonsWithRowSize(2, updateFilmAfterViewedButtons...)
		}).
		AddBack(states.CallUpdateFilmBack).
		Build(session.Lang)
}

func FilmFilters(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	filter := session.GetFilmFiltersByCtx()
	return New().
		AddButtons(getFiltersFilmsButtons(filter, session.Lang)...).
		AddIf(filter.IsEnabled(), func(k *Keyboard) {
			k.AddResetAllFilmsFilters()
		}).
		AddBack(states.CallFilmFiltersBack).
		Build(session.Lang)
}

func FilmsSorting(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	sorting := session.GetFilmSortingByCtx()
	return New().
		AddButtons(getSortingFilmsButtons(sorting, session.Lang)...).
		AddIf(sorting.IsEnabled(), func(k *Keyboard) {
			k.AddResetAllSorting(states.CallFilmSortingAllReset)
		}).
		AddBack(states.CallFilmSortingBack).
		Build(session.Lang)
}

func FilmsNotFound(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddAgain(states.CallFindFilmsAgain).
		AddBack(states.CallFindFilmsBack).
		Build(session.Lang)
}

func NewFilmChangeToken(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddChangeToken().
		AddBack(states.CallFilmsNew).
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
