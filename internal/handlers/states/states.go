package states

const (
	// context keys
	ContextFilm       = "context_film"
	ContextCollection = "context_collection"

	// general states
	CallbackProcessSkip   = "process_skip"
	CallbackProcessCancel = "process_cancel"
	CallbackProcessReset  = "process_reset"

	CallbackIncrease = "increase"
	CallbackDecrease = "decrease"

	CallbackYes = "yes"
	CallbackNo  = "no"

	PrefixSelectStartLang = "select_start_lang_"

	// main menu
	CallbackMainMenu = "main_menu"

	// menu states
	CallbackMenuSelectProfile     = "menu_select_profile"
	CallbackMenuSelectFilms       = "menu_select_films"
	CallbackMenuSelectCollections = "menu_select_collections"
	CallbackMenuSelectSettings    = "menu_select_settings"
	CallbackMenuSelectLogout      = "menu_select_logout"
	CallbackMenuSelectFeedback    = "menu_select_feedback"
	CallbackMenuSelectAdmin       = "menu_select_admin"

	// feedback states
	PrefixFeedbackCategory              = "feedback_category_"
	CallbackFeedbackCategorySuggestions = "feedback_category_suggestions"
	CallbackFeedbackCategoryBugs        = "feedback_category_bugs"
	CallbackFeedbackCategoryIssues      = "feedback_category_issues"
	ProcessFeedbackAwaitingMessage      = "feedback_awaiting_message"

	// logout states
	ProcessLogoutAwaitingConfirm = "logout_awaiting_confirm"

	// settings states
	PrefixSelectLang                           = "select_lang_"
	CallbackSettingsBack                       = "settings_back"
	CallbackSettingsLanguage                   = "settings_language"
	CallbackSettingsKinopoiskToken             = "settings_kinopoisk_token"
	ProcessSettingsAwaitingKinopoiskToken      = "settings_awaiting_kinopoisk_token"
	CallbackSettingsCollectionsPageSize        = "settings_collections_page_size"
	ProcessSettingsCollectionsAwaitingPageSize = "settings_collections_awaiting_page_size"
	CallbackSettingsFilmsPageSize              = "settings_films_page_size"
	ProcessSettingsFilmsAwaitingPageSize       = "settings_films_awaiting_page_size"
	CallbackSettingsObjectsPageSize            = "settings_objects_page_size"
	ProcessSettingsObjectsAwaitingPageSize     = "settings_objects_awaiting_page_size"

	// admin states
	CallbackAdminSelectAdmins    = "admin_select_admins"
	CallbackAdminSelectUsers     = "admin_select_users"
	CallbackAdminSelectBroadcast = "admin_select_broadcast"
	CallbackAdminSelectFeedback  = "admin_select_feedback"

	// admin manage users states
	CallbackAdminManageUsersSelectBack  = "admin_manage_users_select_back"
	CallbackAdminManageUsersSelectFind  = "admin_manage_users_select_find"
	ProcessAdminManageUsersAwaitingFind = "admin_manage_users_awaiting_find"

	// admin users list states
	PrefixAdminUsersListPage        = "admin_users_list_page_"
	CallbackAdminUsersListPagePrev  = "admin_users_list_page_prev"
	CallbackAdminUsersListPageNext  = "admin_users_list_page_next"
	CallbackAdminUsersListPageLast  = "admin_users_list_page_last"
	CallbackAdminUsersListPageFirst = "admin_users_list_page_first"

	// entities states
	PrefixEntities               = "entities_"
	PrefixSelectAdmin            = "select_admin_"
	PrefixSelectAdminUser        = "select_admin_user_"
	PrefixEntitiesListPage       = "entities_page_"
	CallbackEntitiesListBack     = "entities_back"
	CallbackEntitiesListPagePrev = "entities_page_prev"
	CallbackEntitiesListPageNext = "entities_page_next"
	CallbackEntitiesListPageLast = "entities_page_last"
	CallbackEntitiesPageFirst    = "entities_page_first"
	CallbackEntitiesSelectFind   = "entities_select_find"
	ProcessEntitiesAwaitingFind  = "entities_awaiting_find"

	// admin user detail states
	PrefixAdminUserDetail                = "admin_user_detail_"
	PrefixSelectAdminUserRole            = "admin_user_role_select_"
	CallbackAdminUserDetailAgain         = "admin_user_detail_again"
	CallbackAdminUserDetailBack          = "admin_user_detail_back"
	CallbackAdminUserDetailLogs          = "admin_user_detail_logs"
	CallbackAdminUserDetailRole          = "admin_user_detail_role"
	CallbackAdminUserDetailBan           = "admin_user_detail_ban"
	ProcessAdminUserDetailAwaitingReason = "admin_user_detail_awaiting_reason"
	CallbackAdminUserDetailUnban         = "admin_user_detail_unban"
	CallbackAdminUserDetailFeedback      = "admin_user_detail_feedback"

	// admin detail states
	CallbackAdminDetailBack       = "admin_detail_back"
	CallbackAdminDetailAgain      = "admin_detail_again"
	CallbackAdminDetailRaiseRole  = "admin_detail_raise_role"
	CallbackAdminDetailLowerRole  = "admin_detail_lower_role"
	CallbackAdminDetailRemoveRole = "admin_detail_remove_role"

	// admin feedback list states
	PrefixAdminFeedbackList            = "admin_feedback_list_"
	PrefixSelectAdminFeedback          = "select_admin_feedback_"
	PrefixAdminFeedbackListPage        = "admin_feedback_list_page_"
	CallbackAdminFeedbackListBack      = "admin_feedback_list_back"
	CallbackAdminFeedbackListPagePrev  = "admin_feedback_list_page_prev"
	CallbackAdminFeedbackListPageNext  = "admin_feedback_list_page_next"
	CallbackAdminFeedbackListPageLast  = "admin_feedback_list_page_last"
	CallbackAdminFeedbackListPageFirst = "admin_feedback_list_page_first"

	// admin feedback detail states
	CallbackAdminFeedbackDetailBack   = "admin_feedback_detail_back"
	CallbackAdminFeedbackDetailDelete = "admin_feedback_detail_delete"

	// admin user roles states
	CallbackAdminUserRoleSelectUser   = "admin_user_role_select_user"
	CallbackAdminUserRoleSelectHelper = "admin_user_role_select_helper"
	CallbackAdminUserRoleSelectAdmin  = "admin_user_role_select_admin"
	CallbackAdminUserRoleSelectSuper  = "admin_user_role_select_super"

	// admin broadcast states
	CallbackAdminBroadcastSend           = "admin_broadcast_send"
	ProcessAdminBroadcastAwaitingImage   = "admin_broadcast_awaiting_image"
	ProcessAdminBroadcastAwaitingText    = "admin_broadcast_awaiting_text"
	ProcessAdminBroadcastAwaitingPin     = "admin_broadcast_awaiting_pin"
	ProcessAdminBroadcastAwaitingConfirm = "admin_broadcast_awaiting_confirm"

	// profile states
	CallbackProfileSelectUpdate = "profile_select_update"
	CallbackProfileSelectDelete = "profile_select_delete"

	// update profile states
	CallbackUpdateProfileSelectBack      = "update_profile_select_back"
	CallbackUpdateProfileSelectUsername  = "update_profile_select_username"
	ProcessUpdateProfileAwaitingUsername = "update_profile_awaiting_username"
	CallbackUpdateProfileSelectEmail     = "update_profile_select_email"
	ProcessUpdateProfileAwaitingEmail    = "update_profile_awaiting_email"

	// delete profile states
	ProcessDeleteProfileAwaitingConfirm = "delete_profile_awaiting_confirm"

	// films states
	PrefixFilms            = "films_"
	PrefixSelectFilm       = "select_film_"
	PrefixFilmsPage        = "films_page_"
	CallbackFilmsNew       = "films_new"
	CallbackFilmsFilters   = "films_filters"
	CallbackFilmsSorting   = "films_sorting"
	CallbackFilmsManage    = "films_manage"
	CallbackFilmsPageNext  = "films_page_next"
	CallbackFilmsPagePrev  = "films_page_prev"
	CallbackFilmsPageLast  = "films_page_last"
	CallbackFilmsPageFirst = "films_page_first"
	CallbackFilmsBack      = "films_back"
	CallbackFilmsFind      = "films_find"

	// find films states
	PrefixFindFilmsPage           = "find_films_page_"
	CallbackFindFilmsBack         = "find_films_back"
	CallbackFindFilmsPageNext     = "find_films_page_next"
	CallbackFindFilmsPagePrev     = "find_films_page_prev"
	CallbackFindFilmsPageLast     = "find_films_page_last"
	CallbackFindFilmsPageFirst    = "find_films_page_first"
	CallbackFindFilmsAgain        = "find_films_again"
	ProcessFindFilmsAwaitingTitle = "find_films_awaiting_title"

	// find new film states
	PrefixSelectFindNewFilm         = "select_find_new_film_"
	PrefixFindNewFilmPage           = "find_new_film_page_"
	CallbackFindNewFilmBack         = "find_new_film_back"
	ProcessFindNewFilmAwaitingTitle = "find_new_film_awaiting_title"
	CallbackFindNewFilmPageNext     = "find_new_film_page_next"
	CallbackFindNewFilmPagePrev     = "find_new_film_page_prev"
	CallbackFindNewFilmPageLast     = "find_new_film_page_last"
	CallbackFindNewFilmPageFirst    = "find_new_film_page_first"
	CallbackFindNewFilmAgain        = "find_new_film_again"

	// filters films states
	PrefixFiltersFilms                         = "filters_films_"
	PrefixFiltersFilmsSelect                   = "filters_films_select_"
	PrefixFiltersFilmsSelectRange              = "filters_films_select_range_"
	PrefixFiltersFilmsSelectSwitch             = "filters_films_select_switch_"
	PrefixFiltersFilmsAwaiting                 = "filters_films_awaiting_"
	PrefixFiltersFilmsAwaitingRange            = "filters_films_awaiting_range_"
	PrefixFiltersFilmsAwaitingSwitch           = "filters_films_awaiting_switch_"
	CallbackFiltersFilmsBack                   = "filters_films_back"
	CallbackFiltersFilmsAllReset               = "filters_films_all_reset"
	CallbackFiltersFilmsSelectRangeRating      = "filters_films_select_range_rating"
	CallbackFiltersFilmsSelectRangeUserRating  = "filters_films_select_range_user_rating"
	CallbackFiltersFilmsSelectRangeYear        = "filters_films_select_range_year"
	CallbackFiltersFilmsSelectSwitchIsViewed   = "filters_films_select_switch_is_viewed"
	CallbackFiltersFilmsSelectSwitchIsFavorite = "filters_films_select_switch_is_favorite"
	CallbackFiltersFilmsSelectSwitchHasURL     = "filters_films_select_switch_has_url"

	// sorting films states
	PrefixSortingFilms                   = "sorting_films_"
	PrefixSortingFilmsSelect             = "sorting_films_select_"
	PrefixSortingFilmsAwaiting           = "sorting_films_awaiting_"
	CallbackSortingFilmsBack             = "sorting_films_back"
	CallbackSortingFilmsAllReset         = "sorting_films_all_reset"
	ProcessSortingFilmsAwaitingDirection = "sorting_films_awaiting_direction"
	CallbackSortingFilmsSelectTitle      = "sorting_films_select_title"
	CallbackSortingFilmsSelectRating     = "sorting_films_select_rating"
	CallbackSortingFilmsSelectYear       = "sorting_films_select_year"
	CallbackSortingFilmsSelectIsViewed   = "sorting_films_select_is_viewed"
	CallbackSortingFilmsSelectIsFavorite = "sorting_films_select_is_favorite"
	CallbackSortingFilmsSelectUserRating = "sorting_films_select_user_rating"
	CallbackSortingFilmsSelectCreatedAt  = "sorting_films_select_created_at"

	// manage film states
	CallbackManageFilmSelectBack                 = "manage_film_select_back"
	CallbackManageFilmSelectUpdate               = "manage_film_select_update"
	CallbackManageFilmSelectDelete               = "manage_film_select_delete"
	CallbackManageFilmSelectRemoveFromCollection = "manage_film_select_remove_from_collection"

	// new film states
	CallbackNewFilmSelectBack                 = "new_film_select_back"
	CallbackNewFilmSelectManually             = "new_film_select_manually"
	CallbackNewFilmSelectFromURL              = "new_film_select_from_url"
	CallbackNewFilmSelectFind                 = "new_film_select_find"
	CallbackNewFilmSelectChangeKinopoiskToken = "new_film_select_change_kinopoisk_token"
	ProcessNewFilmAwaitingURL                 = "new_film_awaiting_url"
	ProcessNewFilmAwaitingTitle               = "new_film_awaiting_title"
	ProcessNewFilmAwaitingYear                = "new_film_awaiting_year"
	ProcessNewFilmAwaitingGenre               = "new_film_awaiting_genre"
	ProcessNewFilmAwaitingDescription         = "new_film_awaiting_description"
	ProcessNewFilmAwaitingRating              = "new_film_awaiting_rating"
	ProcessNewFilmAwaitingImage               = "new_film_awaiting_image"
	ProcessNewFilmAwaitingComment             = "new_film_awaiting_comment"
	ProcessNewFilmAwaitingViewed              = "new_film_awaiting_viewed"
	ProcessNewFilmAwaitingFilmURL             = "new_film_awaiting_film_url"
	ProcessNewFilmAwaitingUserRating          = "new_film_awaiting_user_rating"
	ProcessNewFilmAwaitingReview              = "new_film_awaiting_review"
	ProcessNewFilmAwaitingKinopoiskToken      = "new_film_awaiting_kinopoisk_token"

	// delete film states
	ProcessDeleteFilmAwaitingConfirm = "delete_film_awaiting_confirm"

	// update film states
	CallbackUpdateFilmSelectBack         = "update_film_select_back"
	CallbackUpdateFilmSelectTitle        = "update_film_select_title"
	ProcessUpdateFilmAwaitingTitle       = "update_film_awaiting_title"
	CallbackUpdateFilmSelectYear         = "update_film_select_year"
	ProcessUpdateFilmAwaitingYear        = "update_film_awaiting_year"
	CallbackUpdateFilmSelectGenre        = "update_film_select_genre"
	ProcessUpdateFilmAwaitingGenre       = "update_film_awaiting_genre"
	CallbackUpdateFilmSelectDescription  = "update_film_select_description"
	ProcessUpdateFilmAwaitingDescription = "update_film_awaiting_description"
	CallbackUpdateFilmSelectRating       = "update_film_select_rating"
	ProcessUpdateFilmAwaitingRating      = "update_film_awaiting_rating"
	CallbackUpdateFilmSelectImage        = "update_film_select_image"
	ProcessUpdateFilmAwaitingImage       = "update_film_awaiting_image"
	CallbackUpdateFilmSelectURL          = "update_film_select_url"
	ProcessUpdateFilmAwaitingURL         = "update_film_awaiting_url"
	CallbackUpdateFilmSelectComment      = "update_film_select_comment"
	ProcessUpdateFilmAwaitingComment     = "update_film_awaiting_comment"
	CallbackUpdateFilmSelectViewed       = "update_film_select_viewed"
	ProcessUpdateFilmAwaitingViewed      = "update_film_awaiting_viewed"
	CallbackUpdateFilmSelectUserRating   = "update_film_select_user_rating"
	ProcessUpdateFilmAwaitingUserRating  = "update_film_awaiting_user_rating"
	CallbackUpdateFilmSelectReview       = "update_film_select_review"
	ProcessUpdateFilmAwaitingReview      = "update_film_awaiting_review"

	// film detail states
	PrefixFilmDetailPage       = "film_detail_page_"
	CallbackFilmDetailPageNext = "film_detail_page_next"
	CallbackFilmDetailPagePrev = "film_detail_page_prev"
	CallbackFilmDetailBack     = "film_detail_back"
	CallbackFilmDetailViewed   = "film_detail_viewed"
	CallbackFilmDetailFavorite = "film_detail_favorite"

	// viewed film states
	ProcessViewedFilmAwaitingUserRating = "viewed_film_awaiting_user_rating"
	ProcessViewedFilmAwaitingReview     = "viewed_film_awaiting_review"

	// collections states
	PrefixCollections                  = "collections_"
	PrefixSelectCollection             = "select_collection_"
	PrefixCollectionsPage              = "collection_page_"
	CallbackCollectionsNew             = "collections_new"
	CallbackCollectionsManage          = "collections_manage"
	CallbackCollectionsFavorite        = "collections_favorite"
	CallbackCollectionsPageNext        = "collections_page_next"
	CallbackCollectionsPagePrev        = "collections_page_prev"
	CallbackCollectionsPageLast        = "collections_page_last"
	CallbackCollectionsPageFirst       = "collections_page_first"
	CallbackCollectionsBack            = "collections_back"
	CallbackCollectionsFind            = "collections_find"
	ProcessFindCollectionsAwaitingName = "find_collections_awaiting_name"
	CallbackCollectionsSorting         = "collections_sorting"

	// sorting collections states
	PrefixSortingCollections                   = "sorting_collections_"
	PrefixSortingCollectionsSelect             = "sorting_collections_select"
	PrefixSortingCollectionsAwaiting           = "sorting_collections_awaiting"
	CallbackSortingCollectionsBack             = "sorting_collections_back"
	CallbackSortingCollectionsAllReset         = "sorting_collections_all_reset"
	ProcessSortingCollectionsAwaitingDirection = "sorting_collections_awaiting_direction"
	CallbackSortingCollectionsSelectIsFavorite = "sorting_collections_select_is_favorite"
	CallbackSortingCollectionsSelectName       = "sorting_collections_select_name"
	CallbackSortingCollectionsSelectCreatedAt  = "sorting_collections_select_created_at"
	CallbackSortingCollectionsSelectTotalFilms = "sorting_collections_select_total_films"

	// find collections states
	PrefixFindCollectionsPage        = "find_collections_page_"
	CallbackFindCollectionsBack      = "find_collections_back"
	CallbackFindCollectionsPageNext  = "find_collections_page_next"
	CallbackFindCollectionsPagePrev  = "find_collections_page_prev"
	CallbackFindCollectionsPageLast  = "find_collections_page_last"
	CallbackFindCollectionsPageFirst = "find_collections_page_first"
	CallbackFindCollectionsAgain     = "find_collections_again"

	// manage collection states
	CallbackManageCollectionSelectBack   = "manage_collection_select_back"
	CallbackManageCollectionSelectUpdate = "manage_collection_select_update"
	CallbackManageCollectionSelectDelete = "manage_collection_select_delete"

	// new collection states
	ProcessNewCollectionAwaitingName        = "new_collection_awaiting_name"
	ProcessNewCollectionAwaitingDescription = "new_collection_awaiting_description"

	// delete collection states
	ProcessDeleteCollectionAwaitingConfirm = "delete_collection_awaiting_confirm"

	// update colllection states
	CallbackUpdateCollectionSelectBack         = "update_collection_select_back"
	CallbackUpdateCollectionSelectName         = "update_collection_select_name"
	ProcessUpdateCollectionAwaitingName        = "update_collection_awaiting_name"
	CallbackUpdateCollectionSelectDescription  = "update_collection_select_description"
	ProcessUpdateCollectionAwaitingDescription = "update_collection_awaiting_description"

	// collection films
	CallbackCollectionFilmsFromFilm       = "collection_films_from_film"
	CallbackCollectionFilmsFromCollection = "collection_films_from_collection"

	// options film to collection states
	CallbackOptionsFilmToCollectionBack     = "options_film_to_collection_back"
	CallbackOptionsFilmToCollectionNew      = "options_film_to_collection_new"
	CallbackOptionsFilmToCollectionExisting = "options_film_to_collection_existing"

	// add collection to film states
	PrefixAddCollectionToFilm              = "add_collection_to_film_"
	PrefixSelectCFCollection               = "select_cf_collection_"
	PrefixAddCollectionToFilmPage          = "add_collection_to_film_page"
	CallbackAddCollectionToFilmBack        = "add_collection_to_film_back"
	CallbackAddCollectionToFilmPagePrev    = "add_collection_to_film_page_prev"
	CallbackAddCollectionToFilmPageNext    = "add_collection_to_film_page_next"
	CallbackAddCollectionToFilmPageLast    = "add_collection_to_film_page_last"
	CallbackAddCollectionToFilmPageFirst   = "add_collection_to_film_page_first"
	CallbackAddCollectionToFilmFind        = "add_collection_to_film_find"
	CallbackAddCollectionToFilmAgain       = "add_collection_to_film_again"
	CallbackAddCollectionToFilmReset       = "add_collection_to_film_reset"
	ProcessAddCollectionToFilmAwaitingName = "add_collection_to_film_awaiting_name"

	// add film to collection states
	PrefixAddFilmToCollection               = "add_film_to_collection_"
	PrefixSelectCFFilm                      = "select_cf_film_"
	PrefixAddFilmToCollectionPage           = "add_film_to_collection_page_"
	CallbackAddFilmToCollectionBack         = "add_film_to_collection_back"
	CallbackAddFilmToCollectionPagePrev     = "add_film_to_collection_page_prev"
	CallbackAddFilmToCollectionPageNext     = "add_film_to_collection_page_next"
	CallbackAddFilmToCollectionPageLast     = "add_film_to_collection_page_last"
	CallbackAddFilmToCollectionPageFirst    = "add_film_to_collection_page_first"
	CallbackAddFilmToCollectionFind         = "add_film_to_collection_find"
	CallbackAddFilmToCollectionAgain        = "add_film_to_collection_again"
	CallbackAddFilmToCollectionReset        = "add_film_to_collection_reset"
	ProcessAddFilmToCollectionAwaitingTitle = "add_film_to_collection_awaiting_title"
)
