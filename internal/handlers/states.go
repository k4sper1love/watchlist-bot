package handlers

import "github.com/k4sper1love/watchlist-bot/internal/models"

const (
	ProcessLogoutStateAwaitingConfirm = "logout_state_awaiting_confirm"

	CallbackSettingsCollectionsPageSize        = "settings_collections_page_size"
	ProcessSettingsCollectionsAwaitingPageSize = "settings_collections_awaiting_page_size"

	CallbackCollectionsNextPage = "collections_next_page"
	CallbackCollectionsPrevPage = "collections_prev_page"

	ProcessNewCollectionAwaitingName        = "new_collection_awaiting_name"
	ProcessNewCollectionAwaitingDescription = "new_collection_awaiting_description"

	CallbackCollectionFilmsNextPage = "collection_films_next_page"
	CallbackCollectionFilmsPrevPage = "collection_films_prev_page"

	CallbackCollectionFilmsDetailNextPage = "collection_films_detail_next_page"
	CallbackCollectionFilmsDetailPrevPage = "collection_films_detail_prev_page"
)

func resetState(session *models.Session) {
	session.State = ""
}

func setState(session *models.Session, state string) {
	session.State = state
}
