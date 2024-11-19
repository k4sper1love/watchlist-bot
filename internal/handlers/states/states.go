package states

const (
	// context Keys
	ContextFilm       = "context_film"
	ContextCollection = "context_collection"

	// general states
	CallbackProcessSkip   = "process_skip"
	CallbackProcessCancel = "process_cancel"

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
	CallbackSettingsCollectionsPageSize        = "settings_collections_page_size"
	ProcessSettingsCollectionsAwaitingPageSize = "settings_collections_awaiting_page_size"
	CallbackSettingsFilmsPageSize              = "settings_films_page_size"
	ProcessSettingsFilmsAwaitingPageSize       = "settings_films_awaiting_page_size"
	CallbackSettingsObjectsPageSize            = "settings_objects_page_size"
	ProcessSettingsObjectsAwaitingPageSize     = "settings_objects_awaiting_page_size"

	// admin states
	CallbackAdminSelectBack                  = "admin_select_back"
	CallbackAdminSelectBackPanel             = "admin_select_back_panel"
	CallbackAdminSelectUserCount             = "admin_select_user_count"
	CallbackAdminSelectBroadcastMessage      = "admin_select_broadcast_message"
	ProcessAdminAwaitingBroadcastMessageText = "admin_awaiting_broadcast_message_text"
	CallbackAdminSelectFeedback              = "admin_select_feedback"
	CallbackAdminSelectUsers                 = "admin_select_users"

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
	CallbackFilmsNew      = "films_new"
	CallbackFilmsManage   = "films_manage"
	CallbackFilmsNextPage = "films_next_page"
	CallbackFilmsPrevPage = "films_prev_page"
	CallbackFilmsBack     = "films_back"

	// manage film states
	CallbackManageFilmSelectBack   = "manage_film_select_back"
	CallbackManageFilmSelectUpdate = "manage_film_select_update"
	CallbackManageFilmSelectDelete = "manage_film_select_delete"

	// new film states
	ProcessNewFilmAwaitingTitle       = "new_film_awaiting_title"
	ProcessNewFilmAwaitingYear        = "new_film_awaiting_year"
	ProcessNewFilmAwaitingGenre       = "new_film_awaiting_genre"
	ProcessNewFilmAwaitingDescription = "new_film_awaiting_description"
	ProcessNewFilmAwaitingRating      = "new_film_awaiting_rating"
	ProcessNewFilmAwaitingImage       = "new_film_awaiting_image"
	ProcessNewFilmAwaitingComment     = "new_film_awaiting_comment"
	ProcessNewFilmAwaitingViewed      = "new_film_awaiting_viewed"
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
