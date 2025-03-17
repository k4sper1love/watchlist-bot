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

// AddNavigation adds navigation buttons (e.g., prev, next, first, last) to the keyboard.
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

// AddCancel adds a cancel button to the keyboard.
func (k *Keyboard) AddCancel() *Keyboard {
	return k.AddButton("", "cancel", states.CallProcessCancel, "", true)
}

// AddSkip adds a skip button to the keyboard.
func (k *Keyboard) AddSkip() *Keyboard {
	return k.AddButton("", "skip", states.CallProcessSkip, "", true)
}

// AddSurvey adds yes/no buttons for surveys or confirmations.
func (k *Keyboard) AddSurvey() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"", "yes", states.CallYes, "", true},
		Button{"", "no", states.CallNo, "", true})
}

// AddBack adds back and main menu buttons to the keyboard.
func (k *Keyboard) AddBack(callbackData string) *Keyboard {
	var buttons []Button

	if callbackData != "" {
		buttons = append(buttons, Button{"←", "back", callbackData, "", true})
	}
	buttons = append(buttons, Button{"🏠 ", "mainMenu", states.CallMainMenu, "", true})

	return k.AddButtonsWithRowSize(len(buttons), buttons...)
}

// AddUpdate adds an update button to the keyboard.
func (k *Keyboard) AddUpdate(callbackData string) *Keyboard {
	return k.AddButton("✏️", "update", callbackData, "", true)
}

// AddDelete adds a delete button to the keyboard.
func (k *Keyboard) AddDelete(callbackData string) *Keyboard {
	return k.AddButton("🗑️", "delete", callbackData, "", true)
}

// AddManage adds a manage button to the keyboard.
func (k *Keyboard) AddManage(callbackData string) *Keyboard {
	return k.AddButton("⚙️", "manage", callbackData, "", true)
}

// AddLanguageSelect adds language selection buttons to the keyboard.
func (k *Keyboard) AddLanguageSelect(languages []string, callback string) *Keyboard {
	var buttons []Button

	for _, lang := range languages {
		buttons = append(buttons, Button{"", fmt.Sprintf(strings.ToUpper(lang)), callback + lang, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

// AddProfileUpdate adds a profile update button to the keyboard.
func (k *Keyboard) AddProfileUpdate() *Keyboard {
	return k.AddButton("✏️", "edit", states.CallProfileUpdate, "", true)
}

// AddProfileDelete adds a profile delete button to the keyboard.
func (k *Keyboard) AddProfileDelete() *Keyboard {
	return k.AddButton("🗑️", "delete", states.CallProfileDelete, "", true)
}

// AddAdminPanel adds an admin panel button to the keyboard.
func (k *Keyboard) AddAdminPanel() *Keyboard {
	return k.AddButton("🛠️", "adminPanel", states.CallMenuAdmin, "", true)
}

// AddUpdateFilmURL adds a button to update a film's URL.
func (k *Keyboard) AddUpdateFilmURL() *Keyboard {
	return k.AddButton("", "filmURL", states.CallUpdateFilmURL, "", true)
}

// AddOpenInBrowser adds a button to open a URL in the browser.
func (k *Keyboard) AddOpenInBrowser(url string) *Keyboard {
	return k.AddButton("", "openInBrowser", "", url, true)
}

// AddFilmViewed adds a button to mark a film as viewed.
func (k *Keyboard) AddFilmViewed() *Keyboard {
	return k.AddButton("✔️", "viewed", states.CallFilmDetailViewed, "", true)
}

// AddFilmSelect adds buttons for selecting films from the current page.
func (k *Keyboard) AddFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
		callback := states.SelectFilm + strconv.Itoa(i)

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), film.Title, callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

// AddFindNewFilmSelect adds buttons for selecting new films.
func (k *Keyboard) AddFindNewFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
		callback := states.SelectNewFilm + strconv.Itoa(i)

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), film.Title, callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

// AddFilmFind adds a button to find a film by title.
func (k *Keyboard) AddFilmFind() *Keyboard {
	return k.AddButton("🔎", "findFilmByTitle", states.CallFilmsFind, "", true)
}

// AddFilmNew adds a button to create a new film.
func (k *Keyboard) AddFilmNew() *Keyboard {
	return k.AddButton("➕", "createFilm", states.CallFilmsNew, "", true)
}

// AddFilmFiltersAndSorting adds buttons for film filters and sorting.
func (k *Keyboard) AddFilmFiltersAndSorting(session *models.Session) *Keyboard {
	filtersEnable := session.GetFilmFiltersByCtx().IsEnabled()
	sortingEnable := session.GetFilmSortingByCtx().IsEnabled()

	return k.AddButtonsWithRowSize(2,
		Button{utils.BoolToEmoji(sortingEnable), "sorting", states.CallFilmsSorting, "", true},
		Button{utils.BoolToEmoji(filtersEnable), "filters", states.CallFilmsFilters, "", true})
}

// AddFilmDelete adds a button to delete a film.
func (k *Keyboard) AddFilmDelete() *Keyboard {
	return k.AddButton("🗑️", "deleteFilm", states.CallManageFilmDelete, "", true)
}

// AddFilmUpdate adds a button to update a film.
func (k *Keyboard) AddFilmUpdate() *Keyboard {
	return k.AddButton("✏️", "updateFilm", states.CallManageFilmUpdate, "", true)
}

// AddFilmManage adds a button to manage a film.
func (k *Keyboard) AddFilmManage() *Keyboard {
	return k.AddButton("⚙️", "manageFilm", states.CallFilmsManage, "", true)
}

// AddFilmRemoveFromCollection adds a button to remove a film from a collection.
func (k *Keyboard) AddFilmRemoveFromCollection() *Keyboard {
	return k.AddButton("🧹󠁝", "removeFromCollection", states.CallManageFilmRemoveFromCollection, "", true)
}

// AddNewFilmManually adds a button to create a new film manually.
func (k *Keyboard) AddNewFilmManually() *Keyboard {
	return k.AddButton("", "manually", states.CallNewFilmManually, "", true)
}

// AddNewFilmFromURL adds a button to create a new film from a URL.
func (k *Keyboard) AddNewFilmFromURL() *Keyboard {
	return k.AddButton("", "fromURL", states.CallNewFilmFromURL, "", true)
}

// AddNewFilmFind adds a button to find a new film.
func (k *Keyboard) AddNewFilmFind() *Keyboard {
	return k.AddButton("", "findFilm", states.CallNewFilmFind, "", true)
}

// AddAgain adds a button to repeat an action.
func (k *Keyboard) AddAgain(callback string) *Keyboard {
	return k.AddButton("↻", "again", callback, "", true)
}

// AddResetAllFilmsFilters adds a button to reset all film filters.
func (k *Keyboard) AddResetAllFilmsFilters() *Keyboard {
	return k.AddButton("🔄", "resetFilters", states.CallFilmFiltersAllReset, "", true)
}

// AddResetFilmsFilter adds a button to reset a specific film filter.
func (k *Keyboard) AddResetFilmsFilter(session *models.Session, filterType string) *Keyboard {
	if session.GetFilmFiltersByCtx().IsFieldEnabled(filterType) {
		return k.AddButton("🔄", "reset", states.CallProcessReset, "", true)
	}
	return k
}

// AddFavorite adds a button to toggle a film's favorite status.
func (k *Keyboard) AddFavorite(isFavorite bool, callback string) *Keyboard {
	messageCode := "makeFavorite"
	if isFavorite {
		messageCode = "removeFavorite"
	}

	return k.AddButton(utils.BoolToStar(!isFavorite), messageCode, callback, "", true)
}

// AddChangeToken adds a button to change a Kinopoisk token.
func (k *Keyboard) AddChangeToken() *Keyboard {
	return k.AddButton("🔄", "changeToken", states.CallNewFilmChangeKinopoiskToken, "", true)
}

// AddCollectionsSelect adds buttons for selecting collections.
func (k *Keyboard) AddCollectionsSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, collection := range session.CollectionsState.Collections {
		itemID := utils.GetItemID(i, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)
		callback := states.SelectCollection + strconv.Itoa(collection.ID)

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), collection.Name, callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

// AddCollectionsNew adds a button to create a new collection.
func (k *Keyboard) AddCollectionsNew() *Keyboard {
	return k.AddButton("➕", "createCollection", states.CallCollectionsNew, "", true)
}

// AddCollectionFiltersAndSorting adds a button for collection sorting.
func (k *Keyboard) AddCollectionFiltersAndSorting(session *models.Session) *Keyboard {
	sortingEnable := session.CollectionsState.Sorting.IsEnabled()
	return k.AddButton(utils.BoolToEmoji(sortingEnable), "sorting", states.CallCollectionsSorting, "", true)
}

// AddSearch adds a search button to the keyboard.
func (k *Keyboard) AddSearch(callback string) *Keyboard {
	return k.AddButton("🔎", "search", callback, "", true)
}

// AddReset adds a reset button to the keyboard.
func (k *Keyboard) AddReset(callback string) *Keyboard {
	return k.AddButton("🔄", "reset", callback, "", true)
}

// AddResetAllSorting adds a button to reset all sorting options.
func (k *Keyboard) AddResetAllSorting(callback string) *Keyboard {
	return k.AddButton("🔄", "resetSorting", callback, "", true)
}

// AddResetSorting adds a button to reset sorting for a specific field.
func (k *Keyboard) AddResetSorting(sorting *models.Sorting) *Keyboard {
	if sorting.IsFieldEnabled(sorting.Field) {
		return k.AddButton("🔄", "reset", states.CallProcessReset, "", true)
	}
	return k
}

// AddSortingDirection adds buttons for sorting direction (increase/decrease order).
func (k *Keyboard) AddSortingDirection() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"⬇️", "decreaseOrder", states.CallDecrease, "", true},
		Button{"⬆️", "increaseOrder", states.CallIncrease, "", true})
}

// AddCollectionFilmFromFilm adds a button to add a collection to a film.
func (k *Keyboard) AddCollectionFilmFromFilm() *Keyboard {
	return k.AddButton("➕", "addCollectionToFilm", states.CallCollectionFilmsFromFilm, "", true)
}

// AddCollectionFilmFromCollection adds a button to add a film to a collection.
func (k *Keyboard) AddCollectionFilmFromCollection() *Keyboard {
	return k.AddButton("➕", "addFilmToCollection", states.CallCollectionFilmsFromCollection, "", true)
}

// AddNewFilmToCollection adds a button to create a new film and add it to a collection.
func (k *Keyboard) AddNewFilmToCollection() *Keyboard {
	return k.AddButton("🆕", "createFilm", states.CallFilmToCollectionOptionNew, "", true)
}

// AddExistingFilmToCollection adds a button to select an existing film and add it to a collection.
func (k *Keyboard) AddExistingFilmToCollection() *Keyboard {
	return k.AddButton("\U0001F7F0", "choiceFromFilms", states.CallFilmToCollectionOptionExisting, "", true)
}

// AddCollectionFilmSelectFilm adds buttons for selecting films to add to a collection.
func (k *Keyboard) AddCollectionFilmSelectFilm(films []apiModels.Film) *Keyboard {
	for _, film := range films {
		k.AddButton(
			utils.BoolToStarOrEmpty(film.IsFavorite),
			fmt.Sprintf("%s (%d)", film.Title, film.ID),
			states.SelectCFFilm+strconv.Itoa(film.ID),
			"",
			false,
		)
	}
	return k
}

// AddCollectionFilmSelectCollection adds buttons for selecting collections to add a film to.
func (k *Keyboard) AddCollectionFilmSelectCollection(collections []apiModels.Collection) *Keyboard {
	for _, collection := range collections {
		k.AddButton(
			utils.BoolToStarOrEmpty(collection.IsFavorite),
			fmt.Sprintf("%s (%d)", collection.Name, collection.ID),
			states.SelectCFCollection+strconv.Itoa(collection.ID),
			"",
			false)
	}
	return k
}

// AddSuperRole adds a button to assign the super admin role.
func (k *Keyboard) AddSuperRole() *Keyboard {
	return k.AddButton("🦸", "superAdmin", states.CallUserDetailRoleSelectSuper, "", true)
}

// AddLogs adds a button to view logs for a user.
func (k *Keyboard) AddLogs() *Keyboard {
	return k.AddButton("💾", "logs", states.CallUserDetailLogs, "", true)
}

// AddManageRole adds a button to manage a user's role.
func (k *Keyboard) AddManageRole() *Keyboard {
	return k.AddButton("🔄", "manageUserRole", states.CallUserDetailRole, "", true)
}

// AddUnban adds a button to unban a user.
func (k *Keyboard) AddUnban() *Keyboard {
	return k.AddButton("✅", "unban", states.CallUserDetailUnban, "", true)
}

// AddBan adds a button to ban a user.
func (k *Keyboard) AddBan() *Keyboard {
	return k.AddButton("❌", "ban", states.CallUserDetailBan, "", true)
}

// AddViewFeedback adds a button to view feedback from a user.
func (k *Keyboard) AddViewFeedback() *Keyboard {
	return k.AddButton("📩", "viewFeedback", states.CallUserDetailFeedback, "", true)
}

// AddRaiseRank adds a button to raise a user's role rank.
func (k *Keyboard) AddRaiseRank() *Keyboard {
	return k.AddButton("⬆️", "raiseRole", states.CallAdminDetailRaiseRole, "", true)
}

// AddLowerRank adds a button to lower a user's role rank.
func (k *Keyboard) AddLowerRank() *Keyboard {
	return k.AddButton("⬇️", "lowerRole", states.CallAdminDetailLowerRole, "", true)
}

// AddRemoveAdminRole adds a button to remove an admin role from a user.
func (k *Keyboard) AddRemoveAdminRole() *Keyboard {
	return k.AddButton("⚠️", "removeAdminRole", states.CallAdminDetailRemoveRole, "", true)
}

// AddFeedbackDelete adds a button to delete a feedback entry.
func (k *Keyboard) AddFeedbackDelete() *Keyboard {
	return k.AddButton("🗑️", "delete", states.CallFeedbackDetailDelete, "", true)
}

// AddSendBroadcast adds a button to send a broadcast message.
func (k *Keyboard) AddSendBroadcast() *Keyboard {
	return k.AddButton("➤", "send", states.CallBroadcastSend, "", true)
}

// AddUserSelect adds buttons for selecting users.
func (k *Keyboard) AddUserSelect(session *models.Session, users []models.Session) *Keyboard {
	buttons := entitySelectButtons(session, users, states.SelectUser)
	return k.AddButtonsWithRowSize(2, buttons...)
}

// AddAdminSelect adds buttons for selecting admins.
func (k *Keyboard) AddAdminSelect(session *models.Session, admins []models.Session) *Keyboard {
	buttons := entitySelectButtons(session, admins, states.SelectAdmin)
	return k.AddButtonsWithRowSize(2, buttons...)
}

// AddFeedbackSelect adds buttons for selecting feedback entries.
func (k *Keyboard) AddFeedbackSelect(session *models.Session, feedbacks []models.Feedback) *Keyboard {
	var buttons []Button

	for i, feedback := range feedbacks {
		itemID := utils.GetItemID(i, session.AdminState.CurrentPage, session.AdminState.PageSize)
		callback := states.SelectFeedback + strconv.Itoa(int(feedback.ID))

		buttons = append(buttons, Button{utils.NumberToEmoji(itemID), "", callback, "", false})
	}

	return k.AddButtonsWithRowSize(2, buttons...)
}

// entitySelectButtons generates buttons for selecting entities (users or admins).
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

// addSortingButton adds a sorting button based on the current sorting state.
func addSortingButton(buttons []Button, sorting *models.Sorting, lang, field, callback string) []Button {
	sortingEnabled := sorting.IsFieldEnabled(field)

	text := translator.Translate(lang, field, nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}

	return append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, callback, "", false})
}

// addFiltersFilmsButton adds a filter button for films based on the current filter state.
func addFiltersFilmsButton(buttons []Button, filter *models.FilmFilters, lang, filterType, callback string, needTranslate bool) []Button {
	filterEnabled := filter.IsFieldEnabled(filterType)

	text := translator.Translate(lang, filterType, nil, nil)
	if filterEnabled {
		value := filter.String(filterType)
		if needTranslate {
			value = translator.Translate(lang, value, nil, nil)
		}
		text += fmt.Sprintf(": %s", value)
	}

	return append(buttons, Button{utils.BoolToEmoji(filterEnabled), text, callback, "", false})
}

// getFiltersFilmsButtons generates buttons for film filters.
func getFiltersFilmsButtons(filter *models.FilmFilters, lang string) []Button {
	var buttons []Button

	buttons = addFiltersFilmsButton(buttons, filter, lang, "is_favorite", states.CallFilmFiltersSelectSwitchIsFavorite, true)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "is_viewed", states.CallFilmFiltersSelectSwitchIsViewed, true)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "year", states.CallFilmFiltersSelectRangeYear, false)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "rating", states.CallFilmFiltersSelectRangeRating, false)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "user_rating", states.CallFilmFiltersSelectRangeUserRating, false)
	buttons = addFiltersFilmsButton(buttons, filter, lang, "has_url", states.CallFilmFiltersSelectSwitchHasURL, true)

	return buttons
}

// getSortingFilmsButtons generates buttons for film sorting options.
func getSortingFilmsButtons(sorting *models.Sorting, lang string) []Button {
	var buttons []Button

	buttons = addSortingButton(buttons, sorting, lang, "is_favorite", states.CallFilmSortingSelectIsFavorite)
	buttons = addSortingButton(buttons, sorting, lang, "is_viewed", states.CallFilmSortingSelectIsViewed)
	buttons = addSortingButton(buttons, sorting, lang, "title", states.CallFilmSortingSelectTitle)
	buttons = addSortingButton(buttons, sorting, lang, "year", states.CallFilmSortingSelectYear)
	buttons = addSortingButton(buttons, sorting, lang, "rating", states.CallFilmSortingSelectRating)
	buttons = addSortingButton(buttons, sorting, lang, "user_rating", states.CallFilmSortingSelectUserRating)
	buttons = addSortingButton(buttons, sorting, lang, "created_at", states.CallFilmSortingSelectCreatedAt)

	return buttons
}

// getSortingCollectionsButtons generates buttons for collection sorting options.
func getSortingCollectionsButtons(sorting *models.Sorting, lang string) []Button {
	var buttons []Button

	buttons = addSortingButton(buttons, sorting, lang, "is_favorite", states.CallCollectionSortingSelectIsFavorite)
	buttons = addSortingButton(buttons, sorting, lang, "title", states.CallCollectionSortingSelectName)
	buttons = addSortingButton(buttons, sorting, lang, "total_films", states.CallCollectionSortingSelectTotalFilms)
	buttons = addSortingButton(buttons, sorting, lang, "created_at", states.CallCollectionSortingSelectCreatedAt)

	return buttons
}
