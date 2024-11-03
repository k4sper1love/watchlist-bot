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

	// new collection states
	ProcessNewCollectionAwaitingName        = "new_collection_awaiting_name"
	ProcessNewCollectionAwaitingDescription = "new_collection_awaiting_description"

	// delete collection states
	ProcessDeleteCollectionAwaitingConfirm = "delete_collection_awaiting_confirm"

	// collection films states
	CallbackCollectionFilmsNew      = "collection_films_new"
	CallbackCollectionFilmsNextPage = "collection_films_next_page"
	CallbackCollectionFilmsPrevPage = "collection_films_prev_page"
	CallbackCollectionFilmsBack     = "collection_films_back"
	CallbackCollectionFilmsDelete   = "collection_films_delete"

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

	// collection film detail states
	CallbackCollectionFilmDetailNextPage = "collection_film_detail_next_page"
	CallbackCollectionFilmDetailPrevPage = "collection_film_detail_prev_page"
	CallbackCollectionFilmDetailBack     = "collection_film_detail_back"
)
