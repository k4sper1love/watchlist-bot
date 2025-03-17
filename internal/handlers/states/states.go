// Package states defines constants for managing user workflows in the Watchlist application.
//
// It includes states and actions for navigation, entity management, settings, admin tasks, and interactions.
package states

const (
	// Сontext
	Ctx           = "context_"         // General context prefix.
	CtxFilm       = Ctx + "film"       // Context for film-related operations.
	CtxCollection = Ctx + "collection" // Context for collection-related operations.

	// General
	Process           = "process_"         // Prefix for general process-related states.
	CallProcessSkip   = Process + "skip"   // Action to skip a process.
	CallProcessCancel = Process + "cancel" // Action to cancel a process.
	CallProcessReset  = Process + "reset"  // Action to reset a process.

	CallIncrease = "increase" // Action to increase a value (e.g., page number).
	CallDecrease = "decrease" // Action to decrease a value (e.g., page number).

	CallYes = "yes" // General confirmation action.
	CallNo  = "no"  // General rejection action.

	SelectStartLang = "select_start_lang_" // Prefix for selecting a language at the start.

	// Main Menu
	CallMainMenu        = "main_menu"          // Action to return to the main menu.
	Menu                = "menu_select_"       // Prefix for main menu options.
	CallMenuProfile     = Menu + "profile"     // Action to navigate to the profile section.
	CallMenuFilms       = Menu + "films"       // Action to navigate to the films section.
	CallMenuCollections = Menu + "collections" // Action to navigate to the collections section.
	CallMenuSettings    = Menu + "settings"    // Action to navigate to the settings section.
	CallMenuLogout      = Menu + "logout"      // Action to log out.
	CallMenuFeedback    = Menu + "feedback"    // Action to submit feedback.
	CallMenuAdmin       = Menu + "admin"       // Action to access admin features.

	// Feedback
	Feedback                        = "feedback_"                      // Prefix for feedback-related states.
	FeedbackAwait                   = Feedback + "await_"              // Prefix for awaiting feedback input.
	FeedbackCategory                = Feedback + "category_"           // Prefix for feedback categories.
	CallFeedbackCategorySuggestions = FeedbackCategory + "suggestions" // Action for suggesting improvements.
	CallFeedbackCategoryBugs        = FeedbackCategory + "bugs"        // Action for reporting bugs.
	CallFeedbackCategoryIssues      = FeedbackCategory + "issues"      // Action for reporting issues.
	AwaitFeedbackMessage            = FeedbackAwait + "message"        // State for awaiting feedback message input.

	// Logout
	Logout             = "logout_"               // Prefix for logout-related states.
	LogoutAwait        = Logout + "await_"       // Prefix for awaiting logout confirmation.
	AwaitLogoutConfirm = LogoutAwait + "confirm" // State for confirming logout.

	// Settings
	SelectLang                       = "select_lang_"                          // Prefix for selecting a language.
	Settings                         = "settings_"                             // Prefix for settings-related states.
	SettingsAwait                    = Settings + "await_"                     // Prefix for awaiting settings input.
	CallSettingsBack                 = Settings + "back"                       // Action to go back from settings.
	CallSettingsLanguage             = Settings + "language"                   // Action to change the language.
	CallSettingsKinopoiskToken       = Settings + "kinopoisk_token"            // Action to update the Kinopoisk API token.
	AwaitSettingsKinopoiskToken      = SettingsAwait + "kinopoisk_token"       // State for awaiting Kinopoisk token input.
	CallSettingsCollectionsPageSize  = Settings + "collections_page_size"      // Action to change collections page size.
	AwaitSettingsCollectionsPageSize = SettingsAwait + "collections_page_size" // State for awaiting collections page size input.
	CallSettingsFilmsPageSize        = Settings + "films_page_size"            // Action to change films page size.
	AwaitSettingsFilmsPageSize       = SettingsAwait + "films_page_size"       // State for awaiting films page size input.
	CallSettingsObjectsPageSize      = Settings + "objects_page_size"          // Action to change objects page size.
	AwaitSettingsObjectsPageSize     = SettingsAwait + "objects_page_size"     // State for awaiting objects page size input.

	// Admin
	Admin              = "admin_select_"     // Prefix for admin-related actions.
	CallAdminAdmins    = Admin + "admins"    // Action to manage admins.
	CallAdminUsers     = Admin + "users"     // Action to manage users.
	CallAdminBroadcast = Admin + "broadcast" // Action to send broadcast messages.
	CallAdminFeedback  = Admin + "feedback"  // Action to view feedback.

	// Entities
	SelectEntity         = "select_entity_"        // Prefix for selecting an entity.
	SelectAdmin          = SelectEntity + "admin_" // Prefix for selecting an admin.
	SelectUser           = SelectEntity + "user_"  // Prefix for selecting a user.
	Entities             = "entities_"             // Prefix for entities-related states.
	EntitiesAwait        = Entities + "await_"     // Prefix for awaiting entities input.
	EntitiesPage         = Entities + "page_"      // Prefix for entities pagination.
	CallEntitiesBack     = Entities + "back"       // Action to go back from entities.
	CallEntitiesPagePrev = EntitiesPage + "prev"   // Action to navigate to the previous entities page.
	CallEntitiesPageNext = EntitiesPage + "next"   // Action to navigate to the next entities page.
	CallEntitiesPageLast = EntitiesPage + "last"   // Action to navigate to the last entities page.
	CallEntitiesFirst    = EntitiesPage + "first"  // Action to navigate to the first entities page.
	CallEntitiesFind     = Entities + "find"       // Action to search for entities.
	AwaitEntitiesFind    = EntitiesAwait + "find"  // State for awaiting entities search input.

	// User Detail
	UserDetail                     = "user_detail_"             // Prefix for user detail-related states.
	UserDetailAwait                = UserDetail + "await_"      // Prefix for awaiting user detail input.
	UserDetailRole                 = UserDetail + "role_"       // Prefix for managing user roles.
	CallUserDetailAgain            = UserDetail + "again"       // Action to reload user details.
	CallUserDetailBack             = UserDetail + "back"        // Action to go back from user details.
	CallUserDetailLogs             = UserDetail + "logs"        // Action to view user logs.
	CallUserDetailBan              = UserDetail + "ban"         // Action to ban a user.
	AwaitUserDetailReason          = UserDetailAwait + "reason" // State for awaiting ban reason input.
	CallUserDetailUnban            = UserDetail + "unban"       // Action to unban a user.
	CallUserDetailFeedback         = UserDetail + "feedback"    // Action to view user feedback.
	CallUserDetailRole             = UserDetail + "role"        // Action to manage user roles.
	CallUserDetailRoleSelectUser   = UserDetailRole + "user"    // Action to set the user role.
	CallUserDetailRoleSelectHelper = UserDetailRole + "helper"  // Action to set the helper role.
	CallUserDetailRoleSelectAdmin  = UserDetailRole + "admin"   // Action to set the admin role.
	CallUserDetailRoleSelectSuper  = UserDetailRole + "super"   // Action to set the super admin role.

	// Admin Detail
	AdminDetail               = "admin_detail_"             // Prefix for admin detail-related states.
	CallAdminDetailBack       = AdminDetail + "back"        // Action to go back from admin details.
	CallAdminDetailAgain      = AdminDetail + "again"       // Action to reload admin details.
	CallAdminDetailRaiseRole  = AdminDetail + "raise_role"  // Action to raise an admin's role.
	CallAdminDetailLowerRole  = AdminDetail + "lower_role"  // Action to lower an admin's role.
	CallAdminDetailRemoveRole = AdminDetail + "remove_role" // Action to remove an admin's role.

	// Feedbacks
	SelectFeedback         = "select_feedback_"      // Prefix for selecting a feedback.
	Feedbacks              = "feedbacks_"            // Prefix for feedbacks-related states.
	FeedbacksPage          = Feedbacks + "page_"     // Prefix for feedbacks pagination.
	CallFeedbacksBack      = Feedbacks + "back"      // Action to go back from feedbacks.
	CallFeedbacksPagePrev  = FeedbacksPage + "prev"  // Action to navigate to the previous feedbacks page.
	CallFeedbacksPageNext  = FeedbacksPage + "next"  // Action to navigate to the next feedbacks page.
	CallFeedbacksPageLast  = FeedbacksPage + "last"  // Action to navigate to the last feedbacks page.
	CallFeedbacksPageFirst = FeedbacksPage + "first" // Action to navigate to the first feedbacks page.

	// Feedback Detail
	FeedbackDetail           = "feedback_detail_"        // Prefix for feedback detail-related states.
	CallFeedbackDetailBack   = FeedbackDetail + "back"   // Action to go back from feedback details.
	CallFeedbackDetailDelete = FeedbackDetail + "delete" // Action to delete a feedback.

	// Broadcast
	Broadcast             = "broadcast_"               // Prefix for broadcast-related states.
	BroadcastAwait        = Broadcast + "await_"       // Prefix for awaiting broadcast input.
	CallBroadcastSend     = Broadcast + "send"         // Action to send a broadcast message.
	AwaitBroadcastImage   = BroadcastAwait + "image"   // State for awaiting broadcast image input.
	AwaitBroadcastText    = BroadcastAwait + "text"    // State for awaiting broadcast text input.
	AwaitBroadcastPin     = BroadcastAwait + "pin"     // State for awaiting broadcast pin input.
	AwaitBroadcastConfirm = BroadcastAwait + "confirm" // State for confirming broadcast.

	// Profile
	Profile           = "profile_"         // Prefix for profile-related states.
	CallProfileUpdate = Profile + "update" // Action to update the profile.
	CallProfileDelete = Profile + "delete" // Action to delete the profile.

	// Update Profile
	UpdateProfile              = "update_profile_"               // Prefix for updating profile-related states.
	UpdateProfileAwait         = UpdateProfile + "await_"        // Prefix for awaiting profile update input.
	CallUpdateProfileBack      = UpdateProfile + "back"          // Action to go back from profile update.
	CallUpdateProfileUsername  = UpdateProfile + "username"      // Action to update the username.
	AwaitUpdateProfileUsername = UpdateProfileAwait + "username" // State for awaiting username input.
	CallUpdateProfileEmail     = UpdateProfile + "email"         // Action to update the email.
	AwaitUpdateProfileEmail    = UpdateProfileAwait + "email"    // State for awaiting email input.

	// Delete Profile
	DeleteProfile             = "delete_profile_"              // Prefix for deleting profile-related states.
	DeleteProfileAwait        = DeleteProfile + "await_"       // Prefix for awaiting profile deletion input.
	AwaitDeleteProfileConfirm = DeleteProfileAwait + "confirm" // State for confirming profile deletion.

	// Films
	SelectFilm         = "select_film_"       // Prefix for selecting a film.
	Films              = "films_"             // Prefix for films-related states.
	FilmsAwait         = Films + "await_"     // Prefix for awaiting films input.
	FilmsPage          = Films + "page_"      // Prefix for films pagination.
	CallFilmsBack      = Films + "back"       // Action to go back from films.
	CallFilmsNew       = Films + "new"        // Action to add a new film.
	CallFilmsFind      = Films + "find"       // Action to search for films.
	CallFilmsFilters   = Films + "filters"    // Action to apply filters to films.
	CallFilmsSorting   = Films + "sorting"    // Action to sort films.
	CallFilmsManage    = Films + "manage"     // Action to manage films.
	CallFilmsPageNext  = FilmsPage + "next"   // Action to navigate to the next films page.
	CallFilmsPagePrev  = FilmsPage + "prev"   // Action to navigate to the previous films page.
	CallFilmsPageLast  = FilmsPage + "last"   // Action to navigate to the last films page.
	CallFilmsPageFirst = FilmsPage + "first"  // Action to navigate to the first films page.
	AwaitFilmsTitle    = FilmsAwait + "title" // State for awaiting film title input.

	// Find Films
	FindFilms              = "find_films_"           // Prefix for finding films-related states.
	FindFilmsPage          = FindFilms + "page_"     // Prefix for paginating found films.
	CallFindFilmsBack      = FindFilms + "back"      // Action to go back from finding films.
	CallFindFilmsAgain     = FindFilms + "again"     // Action to search for films again.
	CallFindFilmsPageNext  = FindFilmsPage + "next"  // Action to navigate to the next found films page.
	CallFindFilmsPagePrev  = FindFilmsPage + "prev"  // Action to navigate to the previous found films page.
	CallFindFilmsPageLast  = FindFilmsPage + "last"  // Action to navigate to the last found films page.
	CallFindFilmsPageFirst = FindFilmsPage + "first" // Action to navigate to the first found films page.

	// Find New Film
	SelectNewFilm            = "select_new_film_"        // Prefix for selecting a new film.
	FindNewFilm              = "find_new_film_"          // Prefix for finding a new film-related states.
	FindNewFilmPage          = FindNewFilm + "page_"     // Prefix for paginating found new films.
	CallFindNewFilmBack      = FindNewFilm + "back"      // Action to go back from finding a new film.
	CallFindNewFilmAgain     = FindNewFilm + "again"     // Action to search for a new film again.
	CallFindNewFilmPageNext  = FindNewFilmPage + "next"  // Action to navigate to the next found new films page.
	CallFindNewFilmPagePrev  = FindNewFilmPage + "prev"  // Action to navigate to the previous found new films page.
	CallFindNewFilmPageLast  = FindNewFilmPage + "last"  // Action to navigate to the last found new films page.
	CallFindNewFilmPageFirst = FindNewFilmPage + "first" // Action to navigate to the first found new films page.

	// Film Filters
	FilmFilters                           = "film_filters_"                         // Prefix for film filters-related states.
	FilmFiltersSelect                     = FilmFilters + "select_"                 // Prefix for selecting film filters.
	FilmFiltersSelectRange                = FilmFiltersSelect + "range_"            // Prefix for range-based film filters.
	FilmFiltersSelectSwitch               = FilmFiltersSelect + "switch_"           // Prefix for switch-based film filters.
	FilmFiltersAwait                      = FilmFilters + "await_"                  // Prefix for awaiting film filter input.
	FilmFiltersAwaitRange                 = FilmFiltersAwait + "range_"             // Prefix for awaiting range-based filter input.
	FilmFiltersAwaitSwitch                = FilmFiltersAwait + "switch_"            // Prefix for awaiting switch-based filter input.
	CallFilmFiltersBack                   = FilmFilters + "back"                    // Action to go back from film filters.
	CallFilmFiltersAllReset               = FilmFilters + "all_reset"               // Action to reset all film filters.
	CallFilmFiltersSelectRangeRating      = FilmFiltersSelectRange + "rating"       // Action to select rating range filter.
	CallFilmFiltersSelectRangeUserRating  = FilmFiltersSelectRange + "user_rating"  // Action to select user rating range filter.
	CallFilmFiltersSelectRangeYear        = FilmFiltersSelectRange + "year"         // Action to select year range filter.
	CallFilmFiltersSelectSwitchIsViewed   = FilmFiltersSelectSwitch + "is_viewed"   // Action to toggle viewed status filter.
	CallFilmFiltersSelectSwitchIsFavorite = FilmFiltersSelectSwitch + "is_favorite" // Action to toggle favorite status filter.
	CallFilmFiltersSelectSwitchHasURL     = FilmFiltersSelectSwitch + "has_url"     // Action to toggle URL availability filter.

	// Film Sorting
	FilmSorting                     = "film_sorting_"                   // Prefix for film sorting-related states.
	FilmSortingSelect               = FilmSorting + "select_"           // Prefix for selecting film sorting options.
	FilmSortingAwait                = FilmSorting + "await_"            // Prefix for awaiting film sorting input.
	CallFilmSortingBack             = FilmSorting + "back"              // Action to go back from film sorting.
	CallFilmSortingAllReset         = FilmSorting + "all_reset"         // Action to reset all film sorting options.
	AwaitFilmSortingDirection       = FilmSortingAwait + "direction"    // State for awaiting sorting direction input.
	CallFilmSortingSelectTitle      = FilmSortingSelect + "title"       // Action to sort films by title.
	CallFilmSortingSelectRating     = FilmSortingSelect + "rating"      // Action to sort films by rating.
	CallFilmSortingSelectYear       = FilmSortingSelect + "year"        // Action to sort films by year.
	CallFilmSortingSelectIsViewed   = FilmSortingSelect + "is_viewed"   // Action to sort films by viewed status.
	CallFilmSortingSelectIsFavorite = FilmSortingSelect + "is_favorite" // Action to sort films by favorite status.
	CallFilmSortingSelectUserRating = FilmSortingSelect + "user_rating" // Action to sort films by user rating.
	CallFilmSortingSelectCreatedAt  = FilmSortingSelect + "created_at"  // Action to sort films by creation date.

	// Manage Film
	ManageFilm                         = "manage_film_"                        // Prefix for managing film-related actions.
	CallManageFilmBack                 = ManageFilm + "back"                   // Action to go back from managing a film.
	CallManageFilmUpdate               = ManageFilm + "update"                 // Action to update a film's details.
	CallManageFilmDelete               = ManageFilm + "delete"                 // Action to delete a film.
	CallManageFilmRemoveFromCollection = ManageFilm + "remove_from_collection" // Action to remove a film from a collection.

	// New Film
	NewFilm                         = "new_film_"                        // Prefix for adding a new film.
	NewFilmAwait                    = NewFilm + "await_"                 // Prefix for awaiting new film input.
	CallNewFilmBack                 = NewFilm + "back"                   // Action to go back from adding a new film.
	CallNewFilmManually             = NewFilm + "manually"               // Action to add a film manually.
	CallNewFilmFromURL              = NewFilm + "from_url"               // Action to add a film from a URL.
	CallNewFilmFind                 = NewFilm + "find"                   // Action to search for a film.
	CallNewFilmChangeKinopoiskToken = NewFilm + "change_kinopoisk_token" // Action to change the Kinopoisk API token.
	AwaitNewFilmFind                = NewFilmAwait + "find"              // State for awaiting film search input.
	AwaitNewFilmFromURL             = NewFilmAwait + "from_url"          // State for awaiting film URL input.
	AwaitNewFilmTitle               = NewFilmAwait + "title"             // State for awaiting film title input.
	AwaitNewFilmYear                = NewFilmAwait + "year"              // State for awaiting film year input.
	AwaitNewFilmGenre               = NewFilmAwait + "genre"             // State for awaiting film genre input.
	AwaitNewFilmDescription         = NewFilmAwait + "description"       // State for awaiting film description input.
	AwaitNewFilmRating              = NewFilmAwait + "rating"            // State for awaiting film rating input.
	AwaitNewFilmImage               = NewFilmAwait + "image"             // State for awaiting film image input.
	AwaitNewFilmComment             = NewFilmAwait + "comment"           // State for awaiting film comment input.
	AwaitNewFilmViewed              = NewFilmAwait + "viewed"            // State for awaiting film viewed status input.
	AwaitNewFilmFilmURL             = NewFilmAwait + "film_url"          // State for awaiting film URL input.
	AwaitNewFilmUserRating          = NewFilmAwait + "user_rating"       // State for awaiting user rating input.
	AwaitNewFilmReview              = NewFilmAwait + "review"            // State for awaiting film review input.
	AwaitNewFilmKinopoiskToken      = NewFilmAwait + "kinopoisk_token"   // State for awaiting Kinopoisk API token input.

	// Delete Film
	DeleteFilm             = "delete_film_"              // Prefix for deleting a film.
	DeleteFilmAwait        = DeleteFilm + "await_"       // Prefix for awaiting film deletion input.
	AwaitDeleteFilmConfirm = DeleteFilmAwait + "confirm" // State for confirming film deletion.

	// Update Film
	UpdateFilm                 = "update_film_"                  // Prefix for updating a film.
	UpdateFilmAwait            = UpdateFilm + "await_"           // Prefix for awaiting film update input.
	CallUpdateFilmBack         = UpdateFilm + "back"             // Action to go back from updating a film.
	CallUpdateFilmTitle        = UpdateFilm + "title"            // Action to update the film title.
	AwaitUpdateFilmTitle       = UpdateFilmAwait + "title"       // State for awaiting film title input.
	CallUpdateFilmYear         = UpdateFilm + "year"             // Action to update the film year.
	AwaitUpdateFilmYear        = UpdateFilmAwait + "year"        // State for awaiting film year input.
	CallUpdateFilmGenre        = UpdateFilm + "genre"            // Action to update the film genre.
	AwaitUpdateFilmGenre       = UpdateFilmAwait + "genre"       // State for awaiting film genre input.
	CallUpdateFilmDescription  = UpdateFilm + "description"      // Action to update the film description.
	AwaitUpdateFilmDescription = UpdateFilmAwait + "description" // State for awaiting film description input.
	CallUpdateFilmRating       = UpdateFilm + "rating"           // Action to update the film rating.
	AwaitUpdateFilmRating      = UpdateFilmAwait + "rating"      // State for awaiting film rating input.
	CallUpdateFilmImage        = UpdateFilm + "image"            // Action to update the film image.
	AwaitUpdateFilmImage       = UpdateFilmAwait + "image"       // State for awaiting film image input.
	CallUpdateFilmURL          = UpdateFilm + "url"              // Action to update the film URL.
	AwaitUpdateFilmURL         = UpdateFilmAwait + "url"         // State for awaiting film URL input.
	CallUpdateFilmComment      = UpdateFilm + "comment"          // Action to update the film comment.
	AwaitUpdateFilmComment     = UpdateFilmAwait + "comment"     // State for awaiting film comment input.
	CallUpdateFilmViewed       = UpdateFilm + "viewed"           // Action to update the film viewed status.
	AwaitUpdateFilmViewed      = UpdateFilmAwait + "viewed"      // State for awaiting film viewed status input.
	CallUpdateFilmUserRating   = UpdateFilm + "user_rating"      // Action to update the user rating.
	AwaitUpdateFilmUserRating  = UpdateFilmAwait + "user_rating" // State for awaiting user rating input.
	CallUpdateFilmReview       = UpdateFilm + "review"           // Action to update the film review.
	AwaitUpdateFilmReview      = UpdateFilmAwait + "review"      // State for awaiting film review input.

	// Film Detail
	FilmDetail             = "film_detail_"          // Prefix for viewing film details.
	FilmDetailPage         = FilmDetail + "page_"    // Prefix for paginating film details.
	CallFilmDetailPageNext = FilmDetailPage + "next" // Action to navigate to the next film detail page.
	CallFilmDetailPagePrev = FilmDetailPage + "prev" // Action to navigate to the previous film detail page.
	CallFilmDetailBack     = FilmDetail + "back"     // Action to go back from film details.
	CallFilmDetailViewed   = FilmDetail + "viewed"   // Action to mark a film as viewed.
	CallFilmDetailFavorite = FilmDetail + "favorite" // Action to mark a film as favorite.

	// Viewed Film
	ViewedFilm                = "viewed_film_"                  // Prefix for viewed film-related states.
	ViewedFilmAwait           = ViewedFilm + "await_"           // Prefix for awaiting viewed film input.
	AwaitViewedFilmUserRating = ViewedFilmAwait + "user_rating" // State for awaiting user rating input.
	AwaitViewedFilmReview     = ViewedFilmAwait + "review"      // State for awaiting film review input.

	// Collections
	SelectCollection         = "select_collection_"       // Prefix for selecting a collection.
	Collections              = "collections_"             // Prefix for collections-related states.
	CollectionsAwait         = Collections + "await_"     // Prefix for awaiting collections input.
	CollectionsPage          = Collections + "page_"      // Prefix for paginating collections.
	CallCollectionsNew       = Collections + "new"        // Action to create a new collection.
	CallCollectionsManage    = Collections + "manage"     // Action to manage a collection.
	CallCollectionsFavorite  = Collections + "favorite"   // Action to mark a collection as favorite.
	CallCollectionsPageNext  = Collections + "page_next"  // Action to navigate to the next collections page.
	CallCollectionsPagePrev  = Collections + "page_prev"  // Action to navigate to the previous collections page.
	CallCollectionsPageLast  = Collections + "page_last"  // Action to navigate to the last collections page.
	CallCollectionsPageFirst = Collections + "page_first" // Action to navigate to the first collections page.
	CallCollectionsBack      = Collections + "back"       // Action to go back from collections.
	CallCollectionsFind      = Collections + "find"       // Action to search for collections.
	CallCollectionsSorting   = Collections + "sorting"    // Action to sort collections.
	AwaitCollectionsName     = CollectionsAwait + "name"  // State for awaiting collection name input.

	// Collection Sorting
	CollectionSorting                     = "collection_sorting_"                   // Prefix for sorting collections.
	CollectionSortingSelect               = CollectionSorting + "select_"           // Prefix for selecting collection sorting options.
	CollectionSortingAwait                = CollectionSorting + "await_"            // Prefix for awaiting collection sorting input.
	CallCollectionSortingBack             = CollectionSorting + "back"              // Action to go back from collection sorting.
	CallCollectionSortingAllReset         = CollectionSorting + "all_reset"         // Action to reset all collection sorting options.
	CallCollectionSortingSelectIsFavorite = CollectionSortingSelect + "is_favorite" // Action to sort by favorite status.
	CallCollectionSortingSelectName       = CollectionSortingSelect + "name"        // Action to sort by collection name.
	CallCollectionSortingSelectCreatedAt  = CollectionSortingSelect + "created_at"  // Action to sort by creation date.
	CallCollectionSortingSelectTotalFilms = CollectionSortingSelect + "total_films" // Action to sort by total films.
	AwaitCollectionSortingDirection       = CollectionSortingAwait + "direction"    // State for awaiting sorting direction input.

	// Find Collections
	FindCollections              = "find_collections_"           // Prefix for finding collections-related states.
	FindCollectionsPage          = FindCollections + "page_"     // Prefix for paginating found collections.
	CallFindCollectionsBack      = FindCollections + "back"      // Action to go back from finding collections.
	CallFindCollectionsPageNext  = FindCollectionsPage + "next"  // Action to navigate to the next found collections page.
	CallFindCollectionsPagePrev  = FindCollectionsPage + "prev"  // Action to navigate to the previous found collections page.
	CallFindCollectionsPageLast  = FindCollectionsPage + "last"  // Action to navigate to the last found collections page.
	CallFindCollectionsPageFirst = FindCollectionsPage + "first" // Action to navigate to the first found collections page.
	CallFindCollectionsAgain     = FindCollections + "again"     // Action to search for collections again.

	// Manage Collection
	ManageCollection           = "manage_collection_"        // Prefix for managing a collection.
	CallManageCollectionBack   = ManageCollection + "back"   // Action to go back from managing a collection.
	CallManageCollectionUpdate = ManageCollection + "update" // Action to update a collection.
	CallManageCollectionDelete = ManageCollection + "delete" // Action to delete a collection.

	// New Collection
	NewCollection                 = "new_collection_"                  // Prefix for creating a new collection.
	NewCollectionAwait            = NewCollection + "await_"           // Prefix for awaiting new collection input.
	AwaitNewCollectionName        = NewCollectionAwait + "name"        // State for awaiting collection name input.
	AwaitNewCollectionDescription = NewCollectionAwait + "description" // State for awaiting collection description input.

	// Delete Collection
	DeleteCollection             = "delete_collection_"              // Prefix for deleting a collection.
	DeleteCollectionAwait        = DeleteCollection + "await_"       // Prefix for awaiting collection deletion input.
	AwaitDeleteCollectionConfirm = DeleteCollectionAwait + "confirm" // State for confirming collection deletion.

	// Update Collection
	UpdateCollection                 = "update_collection_"                  // Prefix for updating a collection.
	UpdateCollectionAwait            = UpdateCollection + "await_"           // Prefix for awaiting collection update input.
	CallUpdateCollectionBack         = UpdateCollection + "back"             // Action to go back from updating a collection.
	CallUpdateCollectionName         = UpdateCollection + "name"             // Action to update the collection name.
	AwaitUpdateCollectionName        = UpdateCollectionAwait + "name"        // State for awaiting collection name input.
	CallUpdateCollectionDescription  = UpdateCollection + "description"      // Action to update the collection description.
	AwaitUpdateCollectionDescription = UpdateCollectionAwait + "description" // State for awaiting collection description input.

	// Collection Films
	CollectionFilmsFrom               = "collection_films_from_"           // Prefix for managing films in a collection.
	CallCollectionFilmsFromFilm       = CollectionFilmsFrom + "film"       // Action to add a film to a collection.
	CallCollectionFilmsFromCollection = CollectionFilmsFrom + "collection" // Action to manage films from another collection.

	// Film to Collection Option
	FilmToCollectionOption             = "film_to_collection_option_"        // Prefix for adding a film to a collection.
	CallFilmToCollectionOptionBack     = FilmToCollectionOption + "back"     // Action to go back from adding a film to a collection.
	CallFilmToCollectionOptionNew      = FilmToCollectionOption + "new"      // Action to create a new collection for the film.
	CallFilmToCollectionOptionExisting = FilmToCollectionOption + "existing" // Action to add the film to an existing collection.

	// Add Collection to Film
	SelectCFCollection               = "select_cf_collection_"           // Prefix for selecting a collection to add to a film.
	AddCollectionToFilm              = "add_collection_to_film_"         // Prefix for adding a collection to a film.
	AddCollectionToFilmPage          = AddCollectionToFilm + "page_"     // Prefix for paginating collections to add to a film.
	AddCollectionToFilmAwait         = AddCollectionToFilm + "await_"    // Prefix for awaiting collection input.
	CallAddCollectionToFilmBack      = AddCollectionToFilm + "back"      // Action to go back from adding a collection to a film.
	CallAddCollectionToFilmPagePrev  = AddCollectionToFilmPage + "prev"  // Action to navigate to the previous collections page.
	CallAddCollectionToFilmPageNext  = AddCollectionToFilmPage + "next"  // Action to navigate to the next collections page.
	CallAddCollectionToFilmPageLast  = AddCollectionToFilmPage + "last"  // Action to navigate to the last collections page.
	CallAddCollectionToFilmPageFirst = AddCollectionToFilmPage + "first" // Action to navigate to the first collections page.
	CallAddCollectionToFilmFind      = AddCollectionToFilm + "find"      // Action to search for collections to add to a film.
	CallAddCollectionToFilmAgain     = AddCollectionToFilm + "again"     // Action to search for collections again.
	CallAddCollectionToFilmReset     = AddCollectionToFilm + "reset"     // Action to reset the collection search.
	AwaitAddCollectionToFilmName     = AddCollectionToFilmAwait + "name" // State for awaiting collection name input.

	// Add Film to Collection
	SelectCFFilm                     = "select_cf_film_"                  // Prefix for selecting a film to add to a collection.
	AddFilmToCollection              = "add_film_to_collection_"          // Prefix for adding a film to a collection.
	AddFilmToCollectionPage          = AddFilmToCollection + "page_"      // Prefix for paginating films to add to a collection.
	AddFilmToCollectionAwait         = AddFilmToCollectionPage + "await_" // Prefix for awaiting film input.
	CallAddFilmToCollectionBack      = AddFilmToCollection + "back"       // Action to go back from adding a film to a collection.
	CallAddFilmToCollectionPagePrev  = AddFilmToCollectionPage + "prev"   // Action to navigate to the previous films page.
	CallAddFilmToCollectionPageNext  = AddFilmToCollectionPage + "next"   // Action to navigate to the next films page.
	CallAddFilmToCollectionPageLast  = AddFilmToCollectionPage + "last"   // Action to navigate to the last films page.
	CallAddFilmToCollectionPageFirst = AddFilmToCollectionPage + "first"  // Action to navigate to the first films page.
	CallAddFilmToCollectionFind      = AddFilmToCollection + "find"       // Action to search for films to add to a collection.
	CallAddFilmToCollectionAgain     = AddFilmToCollection + "again"      // Action to search for films again.
	CallAddFilmToCollectionReset     = AddFilmToCollection + "reset"      // Action to reset the film search.
	AwaitAddFilmToCollectionTitle    = AddFilmToCollectionAwait + "title" // State for awaiting film title input.
)
