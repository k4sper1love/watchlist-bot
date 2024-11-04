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
	CallbackMenuSelectCollections = "menu_select_collections"
	CallbackMenuSelectSettings    = "menu_select_settings"
	CallbackMenuSelectLogout      = "menu_select_logout"

	// settings states
	CallbackSettingsBack                       = "settings_back"
	CallbackSettingsCollectionsPageSize        = "settings_collections_page_size"
	ProcessSettingsCollectionsAwaitingPageSize = "settings_collections_awaiting_page_size"

	// profile states
	CallbackProfileBack = "profile_back"

	// collections states
	CallbackCollectionsNew      = "collections_new"
	CallbackCollectionsNextPage = "collections_next_page"
	CallbackCollectionsPrevPage = "collections_prev_page"
	CallbackCollectionsBack     = "collections_back"
	CallbackCollectionsDelete   = "collections_delete"
	CallbackCollectionsUpdate   = "collections_update"

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
	CallbackCollectionFilmsNextPage = "collection_films_next_page"
	CallbackCollectionFilmsPrevPage = "collection_films_prev_page"
	CallbackCollectionFilmsBack     = "collection_films_back"
	CallbackCollectionFilmsDelete   = "collection_films_delete"
	CallbackCollectionFilmsUpdate   = "collection_films_update"

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
)
