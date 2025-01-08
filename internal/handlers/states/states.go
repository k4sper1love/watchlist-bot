package states

const (
	// context Keys
	ContextFilm       = "context_film"
	ContextCollection = "context_collection"

	// general states
	CallbackProcessSkip   = "process_skip"
	CallbackProcessCancel = "process_cancel"
	CallbackProcessReset  = "process_reset"

	CallbackIncrease  = "increase"
	CallbacktDecrease = "decrease"

	CallbackYes = "yes"
	CallbackNo  = "no"

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
	CallbackFeedbackCategorySuggestions = "feedback_category_suggestions"
	CallbackFeedbackCategoryBugs        = "feedback_category_bugs"
	CallbackFeedbackCategoryOther       = "feedback_category_other"
	ProcessFeedbackAwaitingMessage      = "feedback_awaiting_message"

	// logout states
	ProcessLogoutAwaitingConfirm = "logout_awaiting_confirm"

	// settings states
	CallbackSettingsBack                       = "settings_back"
	CallbackSettingsLanguage                   = "settings_language"
	CallbackSettingsCollectionsPageSize        = "settings_collections_page_size"
	ProcessSettingsCollectionsAwaitingPageSize = "settings_collections_awaiting_page_size"
	CallbackSettingsFilmsPageSize              = "settings_films_page_size"
	ProcessSettingsFilmsAwaitingPageSize       = "settings_films_awaiting_page_size"
	CallbackSettingsObjectsPageSize            = "settings_objects_page_size"
	ProcessSettingsObjectsAwaitingPageSize     = "settings_objects_awaiting_page_size"

	// admin states
	CallbackAdminSelectBackPanel        = "admin_select_back_panel"
	CallbackAdminSelectUserCount        = "admin_select_user_count"
	CallbackAdminSelectBroadcastMessage = "admin_select_broadcast_message"

	// admin states
	CallbackAdminSelectBack      = "admin_select_back"
	CallbackAdminSelectAdmins    = "admin_select_admins"
	CallbackAdminSelectUsers     = "admin_select_users"
	CallbackAdminSelectBroadcast = "admin_select_broadcast"
	CallbackAdminSelectFeedback  = "admin_select_feedback"

	// admin manage users states
	CallbackAdminManageUsersSelectBack  = "admin_manage_users_select_back"
	CallbackAdminManageUsersSelectFind  = "admin_manage_users_select_find"
	ProcessAdminManageUsersAwaitingFind = "admin_manage_users_awaiting_find"

	// admin users list states
	CallbackAdminUsersListPrevPage = "admin_users_list_prev_page"
	CallbackAdminUsersListNextPage = "admin_users_list_next_page"

	// admin list states
	CallbackAdminListBack        = "admin_list_back"
	CallbackAdminListPrevPage    = "admin_list_prev_page"
	CallbackAdminListNextPage    = "admin_list_next_page"
	CallbackAdminListSelectFind  = "admin_list_select_find"
	ProcessAdminListAwaitingFind = "admin_list_awaiting_find"

	// admin user detail states
	CallbackAdminUserDetailBack          = "admin_user_detail_back"
	CallbackAdminUserDetailRole          = "admin_user_detail_role"
	CallbackAdminUserDetailBan           = "admin_user_detail_ban"
	ProcessAdminUserDetailAwaitingReason = "admin_user_detail_awaiting_reason"
	CallbackAdminUserDetailUnban         = "admin_user_detail_unban"
	CallbackAdminUserDetailFeedback      = "admin_user_detail_feedback"

	// admin detail states
	CallbackAdminDetailBack       = "admin_detail_back"
	CallbackAdminDetailRaiseRole  = "admin_detail_raise_role"
	CallbackAdminDetailLowerRole  = "admin_detail_lower_role"
	CallbackAdminDetailRemoveRole = "admin_detail_remove_role"

	// admin feedback list states
	CallbackAdminFeedbackListBack     = "admin_feedback_list_back"
	CallbackAdminFeedbackListPrevPage = "admin_feedback_list_prev_page"
	CallbackAdminFeedbackListNextPage = "admin_feedback_list_next_page"

	// admin feedback detail states
	CallbackAdminFeedbackDetailBack   = "admin_feedback_detail_back"
	CallbackAdminFeedbackDetailDelete = "admin_feedback_detail_delete"

	// admin user roles states
	CallbackAdminUserRoleSelectBack   = "admin_user_role_select_back"
	CallbackAdminUserRoleSelectUser   = "admin_user_role_select_user"
	CallbackAdminUserRoleSelectHelper = "admin_user_role_select_helper"
	CallbackAdminUserRoleSelectAdmin  = "admin_user_role_select_admin"
	CallbackAdminUserRoleSelectSuper  = "admin_user_role_select_super"

	// admin broadcast states
	CallbackAdminBroadcastBack           = "admin_broadcast_back"
	CallbackAdminBroadcastSend           = "admin_broadcast_send"
	ProcessAdminBroadcastAwaitingText    = "admin_broadcast_awaiting_text"
	ProcessAdminBroadcastAwaitingImage   = "admin_broadcast_awaiting_image"
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
	CallbackFilmsNew              = "films_new"
	CallbackFilmsFilters          = "films_filters"
	CallbackFilmsSorting          = "films_sorting"
	CallbackFilmsManage           = "films_manage"
	CallbackFilmsNextPage         = "films_next_page"
	CallbackFilmsPrevPage         = "films_prev_page"
	CallbackFilmsBack             = "films_back"
	CallbackFilmsFind             = "films_find"
	ProcessFindFilmsAwaitingTitle = "find_films_awaiting_title"

	// find films states
	CallbackFindFilmsBack     = "find_films_back"
	CallbackFindFilmsNextPage = "find_films_next_page"
	CallbackFindFilmsPrevPage = "find_films_prev_page"
	CallbackFindFilmsAgain    = "find_films_again"

	// filters films states
	CallbackFiltersFilmsSelectBack       = "filters_films_select_back"
	CallbackFiltersFilmsSelectAllReset   = "filters_films_select_all_reset"
	CallbackFiltersFilmsSelectMinRating  = "filters_films_select_min_rating"
	ProcessFiltersFilmsAwaitingMinRating = "filters_films_awaiting_min_rating"
	CallbackFiltersFilmsSelectMaxRating  = "filters_films_select_max_rating"
	ProcessFiltersFilmsAwaitingMaxRating = "filters_films_awaiting_max_rating"

	// sorting films states
	CallbackSortingFilmsSelectBack       = "sorting_films_select_back"
	CallbackSortingFilmsSelectAllReset   = "sorting_films_select_all_reset"
	ProcessSortingFilmsAwaitingDirection = "sorting_films_awaiting_direction"
	CallbackSortingFilmsSelectID         = "sorting_films_select_id"
	CallbackSortingFilmsSelectTitle      = "sorting_films_select_title"
	CallbackSortingFilmsSelectRating     = "sorting_films_select_rating"

	// manage film states
	CallbackManageFilmSelectBack                 = "manage_film_select_back"
	CallbackManageFilmSelectUpdate               = "manage_film_select_update"
	CallbackManageFilmSelectDelete               = "manage_film_select_delete"
	CallbackManageFilmSelectRemoveFromCollection = "manage_film_select_remove_from_collection"

	// new film states
	CallbackNewFilmSelectBack         = "new_film_select_back"
	CallbackNewFilmSelectManually     = "new_film_select_manually"
	CallbackNewFilmSelectFromURL      = "new_film_select_from_url"
	ProcessNewFilmAwaitingURL         = "new_film_awaiting_url"
	ProcessNewFilmAwaitingTitle       = "new_film_awaiting_title"
	ProcessNewFilmAwaitingYear        = "new_film_awaiting_year"
	ProcessNewFilmAwaitingGenre       = "new_film_awaiting_genre"
	ProcessNewFilmAwaitingDescription = "new_film_awaiting_description"
	ProcessNewFilmAwaitingRating      = "new_film_awaiting_rating"
	ProcessNewFilmAwaitingImage       = "new_film_awaiting_image"
	ProcessNewFilmAwaitingComment     = "new_film_awaiting_comment"
	ProcessNewFilmAwaitingViewed      = "new_film_awaiting_viewed"
	ProcessNewFilmAwaitingFilmURL     = "new_film_awaiting_film_url"
	ProcessNewFilmAwaitingUserRating  = "new_film_awaiting_user_rating"
	ProcessNewFilmAwaitingReview      = "new_film_awaiting_review"

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
	CallbackFilmDetailNextPage = "film_detail_next_page"
	CallbackFilmDetailPrevPage = "film_detail_prev_page"
	CallbackFilmDetailBack     = "film_detail_back"
	CallbackFilmDetailViewed   = "film_detail_viewed"

	// viewed film states
	ProcessViewedFilmAwaitingRating = "viewed_film_awaiting_rating"
	ProcessViewedFilmAwaitingReview = "viewed_film_awaiting_review"

	// collections states
	CallbackCollectionsNew      = "collections_new"
	CallbackCollectionsManage   = "collections_manage"
	CallbackCollectionsNextPage = "collections_next_page"
	CallbackCollectionsPrevPage = "collections_prev_page"
	CallbackCollectionsBack     = "collections_back"

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
	CallbackAddCollectionToFilmBack     = "add_collection_to_film_back"
	CallbackAddCollectionToFilmPrevPage = "add_collection_to_film_prev_page"
	CallbackAddCollectionToFilmNextPage = "add_collection_to_film_next_page"

	// add film to collection states
	CallbackAddFilmToCollectionBack     = "add_film_to_collection_back"
	CallbackAddFilmToCollectionPrevPage = "add_film_to_collection_prev_page"
	CallbackAddFilmToCollectionNextPage = "add_film_to_collection_next_page"
)
