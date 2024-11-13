package states

const (
	// general states
	CallbackProcessSkip   = "process_skip"
	CallbackProcessCancel = "process_cancel"

	CallbackYes = "yes"
	CallbackNo  = "no"

	// logout states
	ProcessLogoutAwaitingConfirm = "logout_awaiting_confirm"

	// menu states
	CallbackMenuSelectProfile     = "menu_select_profile"
	CallbackMenuSelectFilms       = "menu_select_films"
	CallbackMenuSelectCollections = "menu_select_collections"
	CallbackMenuSelectSettings    = "menu_select_settings"
	CallbackMenuSelectLogout      = "menu_select_logout"
	CallbackMenuSelectAdmin       = "menu_select_admin"

	// settings states
	CallbackSettingsBack                       = "settings_back"
	CallbackSettingsCollectionsPageSize        = "settings_collections_page_size"
	ProcessSettingsCollectionsAwaitingPageSize = "settings_collections_awaiting_page_size"

	// admin states
	CallbackAdminSelectBack                  = "admin_select_back"
	CallbackAdminSelectUserCount             = "admin_select_user_count"
	CallbackAdminSelectBroadcastMessage      = "admin_select_broadcast_message"
	ProcessAdminAwaitingBroadcastMessageText = "admin_awaiting_broadcast_message_text"

	// profile states
	CallbackProfileSelectUpdate = "profile_select_update"
	CallbackProfileSelectDelete = "profile_select_delete"
	CallbackProfileSelectBack   = "profile_select_back"

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
	CallbackViewedFilmBack          = "viewed_film_back"
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

	// collection films states
	CallbackCollectionFilmsNew      = "collection_films_new"
	CallbackCollectionFilmsManage   = "collection_films_manage"
	CallbackCollectionFilmsNextPage = "collection_films_next_page"
	CallbackCollectionFilmsPrevPage = "collection_films_prev_page"
	CallbackCollectionFilmsBack     = "collection_films_back"

	// manage collection film states
	CallbackManageCollectionFilmSelectBack   = "manage_collection_film_select_back"
	CallbackManageCollectionFilmSelectUpdate = "manage_collection_film_select_update"
	CallbackManageCollectionFilmSelectDelete = "manage_collection_film_select_delete"

	// new collection film states
	ProcessNewCollectionFilmAwaitingTitle       = "new_collection_film_awaiting_title"
	ProcessNewCollectionFilmAwaitingYear        = "new_collection_film_awaiting_year"
	ProcessNewCollectionFilmAwaitingGenre       = "new_collection_film_awaiting_genre"
	ProcessNewCollectionFilmAwaitingDescription = "new_collection_film_awaiting_description"
	ProcessNewCollectionFilmAwaitingRating      = "new_collection_film_awaiting_rating"
	ProcessNewCollectionFilmAwaitingImage       = "new_collection_film_awaiting_image"
	ProcessNewCollectionFilmAwaitingComment     = "new_collection_film_awaiting_comment"
	ProcessNewCollectionFilmAwaitingViewed      = "new_collection_film_awaiting_viewed"
	ProcessNewCollectionFilmAwaitingUserRating  = "new_collection_film_awaiting_user_rating"
	ProcessNewCollectionFilmAwaitingReview      = "new_collection_film_awaiting_review"

	// delete collection film states
	ProcessDeleteCollectionFilmAwaitingConfirm = "delete_collection_film_awaiting_confirm"

	// update collection film states
	CallbackUpdateCollectionFilmSelectBack         = "update_collection_film_select_back"
	CallbackUpdateCollectionFilmSelectTitle        = "update_collection_film_select_title"
	ProcessUpdateCollectionFilmAwaitingTitle       = "update_collection_film_awaiting_title"
	CallbackUpdateCollectionFilmSelectYear         = "update_collection_film_select_year"
	ProcessUpdateCollectionFilmAwaitingYear        = "update_collection_film_awaiting_year"
	CallbackUpdateCollectionFilmSelectGenre        = "update_collection_film_select_genre"
	ProcessUpdateCollectionFilmAwaitingGenre       = "update_collection_film_awaiting_genre"
	CallbackUpdateCollectionFilmSelectDescription  = "update_collection_film_select_description"
	ProcessUpdateCollectionFilmAwaitingDescription = "update_collection_film_awaiting_description"
	CallbackUpdateCollectionFilmSelectRating       = "update_collection_film_select_rating"
	ProcessUpdateCollectionFilmAwaitingRating      = "update_collection_film_awaiting_rating"
	CallbackUpdateCollectionFilmSelectImage        = "update_collection_film_select_image"
	ProcessUpdateCollectionFilmAwaitingImage       = "update_collection_film_awaiting_image"
	CallbackUpdateCollectionFilmSelectComment      = "update_collection_film_select_comment"
	ProcessUpdateCollectionFilmAwaitingComment     = "update_collection_film_awaiting_comment"
	CallbackUpdateCollectionFilmSelectViewed       = "update_collection_film_select_viewed"
	ProcessUpdateCollectionFilmAwaitingViewed      = "update_collection_film_awaiting_viewed"
	CallbackUpdateCollectionFilmSelectUserRating   = "update_collection_film_select_user_rating"
	ProcessUpdateCollectionFilmAwaitingUserRating  = "update_collection_film_awaiting_user_rating"
	CallbackUpdateCollectionFilmSelectReview       = "update_collection_film_select_review"
	ProcessUpdateCollectionFilmAwaitingReview      = "update_collection_film_awaiting_review"

	// collection film detail states
	CallbackCollectionFilmDetailNextPage = "collection_film_detail_next_page"
	CallbackCollectionFilmDetailPrevPage = "collection_film_detail_prev_page"
	CallbackCollectionFilmDetailBack     = "collection_film_detail_back"
	CallbackCollectionFilmDetailViewed   = "collection_film_detail_viewed"

	// viewed collection film states
	CallbackViewedCollectionFilmBack          = "viewed_collection_film_back"
	ProcessViewedCollectionFilmAwaitingRating = "viewed_collection_film_awaiting_rating"
	ProcessViewedCollectionFilmAwaitingReview = "viewed_collection_film_awaiting_review"
)
