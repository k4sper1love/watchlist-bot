package keyboards

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strconv"
	"strings"
)

func (k *Keyboard) AddNavigation(currentPage, lastPage int, prefixCallback string, needAdditional bool) *Keyboard {
	var buttons []Button

	if currentPage > 1 {
		if needAdditional {
			buttons = append(buttons, Button{Emoji: "⏮", Callback: prefixCallback + "first"})
		}
		buttons = append(buttons, Button{Emoji: "⬅", Callback: prefixCallback + "prev"})
	}

	if currentPage < lastPage {
		buttons = append(buttons, Button{Emoji: "➡", Callback: prefixCallback + "next"})
		if needAdditional {
			buttons = append(buttons, Button{Emoji: "⏭", Text: "", Callback: prefixCallback + "last"})
		}
	}

	if len(buttons) > 0 {
		return k.AddRow(buttons...)
	}
	return k
}

func (k *Keyboard) AddCancel() *Keyboard {
	return k.AddButton("", "cancel", states.CallbackProcessCancel, "", true)
}

func (k *Keyboard) AddSkip() *Keyboard {
	return k.AddButton("", "skip", states.CallbackProcessSkip, "", true)
}

func (k *Keyboard) AddSurvey() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"", "yes", states.CallbackYes, "", true},
		Button{"", "no", states.CallbackNo, "", true})
}

func (k *Keyboard) AddBack(callbackData string) *Keyboard {
	var buttons []Button

	if callbackData != "" {
		buttons = append(buttons, Button{"←", "back", callbackData, "", true})
	}
	buttons = append(buttons, Button{"🏠 ", "mainMenu", states.CallbackMainMenu, "", true})

	return k.AddButtonsWithRowSize(len(buttons), buttons...)
}

func (k *Keyboard) AddUpdate(callbackData string) *Keyboard {
	return k.AddButton("✏️", "update", callbackData, "", true)
}

func (k *Keyboard) AddDelete(callbackData string) *Keyboard {
	return k.AddButton("🗑️", "delete", callbackData, "", true)
}

func (k *Keyboard) AddManage(callbackData string) *Keyboard {
	return k.AddButton("⚙️", "manage", callbackData, "", true)
}

func (k *Keyboard) AddLanguageSelect(languages []string, callback string) *Keyboard {
	var buttons []Button

	for _, lang := range languages {
		buttons = append(buttons, Button{"", fmt.Sprintf(strings.ToUpper(lang)), callback + lang, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

func (k *Keyboard) AddProfileUpdate() *Keyboard {
	return k.AddButton("✏️", "edit", states.CallbackProfileSelectUpdate, "", true)
}

func (k *Keyboard) AddProfileDelete() *Keyboard {
	return k.AddButton("🗑️", "delete", states.CallbackProfileSelectDelete, "", true)
}

func (k *Keyboard) AddAdminPanel() *Keyboard {
	return k.AddButton("🛠️", "adminPanel", states.CallbackMenuSelectAdmin, "", true)
}

func (k *Keyboard) AddUpdateFilmURL() *Keyboard {
	return k.AddButton("", "filmURL", states.CallbackUpdateFilmSelectURL, "", true)
}

func (k *Keyboard) AddOpenInBrowser(url string) *Keyboard {
	return k.AddButton("", "openInBrowser", "", url, true)
}

func (k *Keyboard) AddFilmViewed() *Keyboard {
	return k.AddButton("✔️", "viewed", states.CallbackFilmDetailViewed, "", true)
}

func (k *Keyboard) AddFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
		callback := states.PrefixSelectFilm + strconv.Itoa(i)

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), film.Title, callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

func (k *Keyboard) AddFindNewFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
		callback := states.PrefixSelectFindNewFilm + strconv.Itoa(i)

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), film.Title, callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

func (k *Keyboard) AddFilmFind() *Keyboard {
	return k.AddButton("🔎", "findFilmByTitle", states.CallbackFilmsFind, "", true)
}

func (k *Keyboard) AddFilmNew() *Keyboard {
	return k.AddButton("➕", "createFilm", states.CallbackFilmsNew, "", true)
}

func (k *Keyboard) AddFilmFiltersAndSorting(session *models.Session) *Keyboard {
	filtersEnable := session.GetFilmsFiltersByContext().IsFiltersEnabled()
	sortingEnable := session.GetFilmsSortingByContext().IsSortingEnabled()

	return k.AddButtonsWithRowSize(2,
		Button{utils.BoolToEmoji(sortingEnable), "sorting", states.CallbackFilmsSorting, "", true},
		Button{utils.BoolToEmoji(filtersEnable), "filters", states.CallbackFilmsFilters, "", true})
}

func (k *Keyboard) AddFilmDelete() *Keyboard {
	return k.AddButton("🗑️", "deleteFilm", states.CallbackManageFilmSelectDelete, "", true)
}

func (k *Keyboard) AddFilmUpdate() *Keyboard {
	return k.AddButton("✏️", "updateFilm", states.CallbackManageFilmSelectUpdate, "", true)
}

func (k *Keyboard) AddFilmManage() *Keyboard {
	return k.AddButton("⚙️", "manageFilm", states.CallbackFilmsManage, "", true)
}

func (k *Keyboard) AddFilmRemoveFromCollection() *Keyboard {
	return k.AddButton("🧹󠁝", "removeFromCollection", states.CallbackManageFilmSelectRemoveFromCollection, "", true)
}

func (k *Keyboard) AddNewFilmManually() *Keyboard {
	return k.AddButton("", "manually", states.CallbackNewFilmSelectManually, "", true)
}

func (k *Keyboard) AddNewFilmFromURL() *Keyboard {
	return k.AddButton("", "fromURL", states.CallbackNewFilmSelectFromURL, "", true)
}

func (k *Keyboard) AddNewFilmFind() *Keyboard {
	return k.AddButton("", "findFilm", states.CallbackNewFilmSelectFind, "", true)
}

func (k *Keyboard) AddAgain(callback string) *Keyboard {
	return k.AddButton("↻", "again", callback, "", true)
}

func (k *Keyboard) AddResetAllFilmsFilters() *Keyboard {
	return k.AddButton("🔄", "resetFilters", states.CallbackFiltersFilmsAllReset, "", true)
}

func (k *Keyboard) AddResetFilmsFilter(session *models.Session, filterType string) *Keyboard {
	if session.GetFilmsFiltersByContext().IsFilterEnabled(filterType) {
		return k.AddButton("🔄", "reset", states.CallbackProcessReset, "", true)
	}
	return k
}

func (k *Keyboard) AddFavorite(isFavorite bool, callback string) *Keyboard {
	messageCode := "makeFavorite"
	if isFavorite {
		messageCode = "removeFavorite"
	}

	return k.AddButton(utils.BoolToStar(!isFavorite), messageCode, callback, "", true)
}

func (k *Keyboard) AddChangeToken() *Keyboard {
	return k.AddButton("🔄", "changeToken", states.CallbackNewFilmSelectChangeKinopoiskToken, "", true)
}

func (k *Keyboard) AddCollectionsSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, collection := range session.CollectionsState.Collections {
		itemID := utils.GetItemID(i, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)
		callback := states.PrefixSelectCollection + strconv.Itoa(collection.ID)

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), collection.Name, callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

func (k *Keyboard) AddCollectionsNew() *Keyboard {
	return k.AddButton("➕", "createCollection", states.CallbackCollectionsNew, "", true)
}

func (k *Keyboard) AddCollectionFiltersAndSorting(session *models.Session) *Keyboard {
	sortingEnable := session.CollectionsState.Sorting.IsSortingEnabled()
	return k.AddButton(utils.BoolToEmoji(sortingEnable), "sorting", states.CallbackCollectionsSorting, "", true)
}

func (k *Keyboard) AddSearch(callback string) *Keyboard {
	return k.AddButton("🔎", "search", callback, "", true)
}

func (k *Keyboard) AddReset(callback string) *Keyboard {
	return k.AddButton("🔄", "reset", callback, "", true)
}

func (k *Keyboard) AddResetAllSorting(callback string) *Keyboard {
	return k.AddButton("🔄", "resetSorting", callback, "", true)
}

func (k *Keyboard) AddResetSorting(sorting *models.Sorting) *Keyboard {
	if sorting.IsSortingFieldEnabled(sorting.Field) {
		return k.AddButton("🔄", "reset", states.CallbackProcessReset, "", true)
	}
	return k
}

func (k *Keyboard) AddSortingDirection() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"⬇️", "decreaseOrder", states.CallbackDecrease, "", true},
		Button{"⬆️", "increaseOrder", states.CallbackIncrease, "", true})
}

func (k *Keyboard) AddCollectionFilmFromFilm() *Keyboard {
	return k.AddButton("➕", "addCollectionToFilm", states.CallbackCollectionFilmsFromFilm, "", true)
}

func (k *Keyboard) AddCollectionFilmFromCollection() *Keyboard {
	return k.AddButton("➕", "addFilmToCollection", states.CallbackCollectionFilmsFromCollection, "", true)
}

func (k *Keyboard) AddNewFilmToCollection() *Keyboard {
	return k.AddButton("🆕", "createFilm", states.CallbackOptionsFilmToCollectionNew, "", true)
}

func (k *Keyboard) AddExistingFilmToCollection() *Keyboard {
	return k.AddButton("\U0001F7F0", "choiceFromFilms", states.CallbackOptionsFilmToCollectionExisting, "", true)
}

func (k *Keyboard) AddCollectionFilmSelectFilm(films []apiModels.Film) *Keyboard {
	for _, film := range films {
		k.AddButton(
			utils.BoolToStarOrEmpty(film.IsFavorite),
			fmt.Sprintf("%s (%d)", film.Title, film.ID),
			states.PrefixSelectCFFilm+strconv.Itoa(film.ID),
			"",
			false,
		)
	}
	return k
}

func (k *Keyboard) AddCollectionFilmSelectCollection(collections []apiModels.Collection) *Keyboard {
	for _, collection := range collections {
		k.AddButton(
			utils.BoolToStarOrEmpty(collection.IsFavorite),
			fmt.Sprintf("%s (%d)", collection.Name, collection.ID),
			states.PrefixSelectCFCollection+strconv.Itoa(collection.ID),
			"",
			false)
	}
	return k
}

func (k *Keyboard) AddSuperRole() *Keyboard {
	return k.AddButton("🦸", "superAdmin", states.CallbackAdminUserRoleSelectSuper, "", true)
}

func (k *Keyboard) AddLogs() *Keyboard {
	return k.AddButton("💾", "logs", states.CallbackAdminUserDetailLogs, "", true)
}

func (k *Keyboard) AddManageRole() *Keyboard {
	return k.AddButton("🔄", "manageUserRole", states.CallbackAdminUserDetailRole, "", true)
}

func (k *Keyboard) AddUnban() *Keyboard {
	return k.AddButton("✅", "unban", states.CallbackAdminUserDetailUnban, "", true)
}

func (k *Keyboard) AddBan() *Keyboard {
	return k.AddButton("❌", "ban", states.CallbackAdminUserDetailBan, "", true)
}

func (k *Keyboard) AddViewFeedback() *Keyboard {
	return k.AddButton("📩", "viewFeedback", states.CallbackAdminUserDetailFeedback, "", true)
}

func (k *Keyboard) AddRaiseRank() *Keyboard {
	return k.AddButton("⬆️", "raiseRole", states.CallbackAdminDetailRaiseRole, "", true)
}

func (k *Keyboard) AddLowerRank() *Keyboard {
	return k.AddButton("⬇️", "lowerRole", states.CallbackAdminDetailLowerRole, "", true)
}

func (k *Keyboard) AddRemoveAdminRole() *Keyboard {
	return k.AddButton("⚠️", "removeAdminRole", states.CallbackAdminDetailRemoveRole, "", true)
}

func (k *Keyboard) AddFeedbackDelete() *Keyboard {
	return k.AddButton("🗑️", "delete", states.CallbackAdminFeedbackDetailDelete, "", true)
}

func (k *Keyboard) AddSendBroadcast() *Keyboard {
	return k.AddButton("➤", "send", states.CallbackAdminBroadcastSend, "", true)
}

func (k *Keyboard) AddUserSelect(session *models.Session, users []models.Session) *Keyboard {
	buttons := entitySelectButtons(session, users, states.PrefixSelectAdminUser)
	return k.AddButtonsWithRowSize(2, buttons...)
}

func (k *Keyboard) AddAdminSelect(session *models.Session, admins []models.Session) *Keyboard {
	buttons := entitySelectButtons(session, admins, states.PrefixSelectAdmin)
	return k.AddButtonsWithRowSize(2, buttons...)
}

func (k *Keyboard) AddFeedbackSelect(session *models.Session, feedbacks []models.Feedback) *Keyboard {
	var buttons []Button

	for i, feedback := range feedbacks {
		itemID := utils.GetItemID(i, session.AdminState.CurrentPage, session.AdminState.PageSize)
		callback := states.PrefixSelectAdminFeedback + strconv.Itoa(int(feedback.ID))

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), "", callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

func entitySelectButtons(session *models.Session, users []models.Session, rawCallback string) []Button {
	var buttons []Button

	for i, user := range users {
		itemID := utils.GetItemID(i, session.AdminState.CurrentPage, session.AdminState.PageSize)
		callback := rawCallback + strconv.Itoa(user.TelegramID)
		text := strconv.Itoa(user.TelegramID)
		if user.TelegramUsername != "" {
			text += fmt.Sprintf(" (@%s)", user.TelegramUsername)
		}

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), text, callback, "", false})
	}

	return buttons
}

func addSortingButton(buttons []Button, sorting *models.Sorting, lang, field, callback string) []Button {
	sortingEnabled := sorting.IsSortingFieldEnabled(field)

	text := translator.Translate(lang, field, nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}

	return append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, callback, "", false})
}

func addFiltersFilmsButton(buttons []Button, filter *models.FiltersFilm, lang, filterType, callback string, needTranslate bool) []Button {
	filterEnabled := filter.IsFilterEnabled(filterType)

	text := translator.Translate(lang, filterType, nil, nil)
	if filterEnabled {
		value := filter.ValueToString(filterType)
		if needTranslate {
			value = translator.Translate(lang, value, nil, nil)
		}
		text += fmt.Sprintf(": %s", value)
	}

	return append(buttons, Button{utils.BoolToEmoji(filterEnabled), text, callback, "", false})
}

func getFiltersFilmsButtons(filter *models.FiltersFilm, lang string) []Button {
	var buttons []Button

	buttons = addFiltersFilmsButton(buttons, filter, lang, "is_favorite", states.CallbackFiltersFilmsSelectSwitchIsFavorite, true)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "is_viewed", states.CallbackFiltersFilmsSelectSwitchIsViewed, true)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "year", states.CallbackFiltersFilmsSelectRangeYear, false)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "rating", states.CallbackFiltersFilmsSelectRangeRating, false)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "user_rating", states.CallbackFiltersFilmsSelectRangeUserRating, false)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "has_url", states.CallbackFiltersFilmsSelectSwitchHasURL, true)

	return buttons
}

func getSortingFilmsButtons(sorting *models.Sorting, lang string) []Button {
	var buttons []Button

	buttons = addSortingButton(buttons, sorting, lang, "is_favorite", states.CallbackSortingFilmsSelectIsFavorite)
	buttons = addSortingButton(buttons, sorting, lang, "is_viewed", states.CallbackSortingFilmsSelectIsViewed)
	buttons = addSortingButton(buttons, sorting, lang, "title", states.CallbackSortingFilmsSelectTitle)
	buttons = addSortingButton(buttons, sorting, lang, "year", states.CallbackSortingFilmsSelectYear)
	buttons = addSortingButton(buttons, sorting, lang, "rating", states.CallbackSortingFilmsSelectRating)
	buttons = addSortingButton(buttons, sorting, lang, "user_rating", states.CallbackSortingFilmsSelectUserRating)
	buttons = addSortingButton(buttons, sorting, lang, "created_at", states.CallbackSortingFilmsSelectCreatedAt)

	return buttons
}

func getSortingCollectionsButtons(sorting *models.Sorting, lang string) []Button {
	var buttons []Button

	buttons = addSortingButton(buttons, sorting, lang, "is_favorite", states.CallbackSortingCollectionsSelectIsFavorite)
	buttons = addSortingButton(buttons, sorting, lang, "title", states.CallbackSortingCollectionsSelectName)
	buttons = addSortingButton(buttons, sorting, lang, "total_films", states.CallbackSortingCollectionsSelectTotalFilms)
	buttons = addSortingButton(buttons, sorting, lang, "created_at", states.CallbackSortingCollectionsSelectCreatedAt)

	return buttons
}
