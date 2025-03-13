package states

const (
	// context
	Ctx           = "context_"
	CtxFilm       = Ctx + "film"
	CtxCollection = Ctx + "collection"

	// general
	Process           = "process_"
	CallProcessSkip   = Process + "skip"
	CallProcessCancel = Process + "cancel"
	CallProcessReset  = Process + "reset"

	CallIncrease = "increase"
	CallDecrease = "decrease"

	CallYes = "yes"
	CallNo  = "no"

	SelectStartLang = "select_start_lang_"

	// main menu
	CallMainMenu        = "main_menu"
	Menu                = "menu_select_"
	CallMenuProfile     = Menu + "profile"
	CallMenuFilms       = Menu + "films"
	CallMenuCollections = Menu + "collections"
	CallMenuSettings    = Menu + "settings"
	CallMenuLogout      = Menu + "logout"
	CallMenuFeedback    = Menu + "feedback"
	CallMenuAdmin       = Menu + "admin"

	// feedback
	Feedback                        = "feedback_"
	FeedbackAwait                   = Feedback + "await_"
	FeedbackCategory                = Feedback + "category_"
	CallFeedbackCategorySuggestions = FeedbackCategory + "suggestions"
	CallFeedbackCategoryBugs        = FeedbackCategory + "bugs"
	CallFeedbackCategoryIssues      = FeedbackCategory + "issues"
	AwaitFeedbackMessage            = FeedbackAwait + "message"

	// logout
	Logout             = "logout_"
	LogoutAwait        = Logout + "await_"
	AwaitLogoutConfirm = LogoutAwait + "confirm"

	// settings
	SelectLang                       = "select_lang_"
	Settings                         = "settings_"
	SettingsAwait                    = Settings + "await_"
	CallSettingsBack                 = Settings + "back"
	CallSettingsLanguage             = Settings + "language"
	CallSettingsKinopoiskToken       = Settings + "kinopoisk_token"
	AwaitSettingsKinopoiskToken      = SettingsAwait + "kinopoisk_token"
	CallSettingsCollectionsPageSize  = Settings + "collections_page_size"
	AwaitSettingsCollectionsPageSize = SettingsAwait + "collections_page_size"
	CallSettingsFilmsPageSize        = Settings + "films_page_size"
	AwaitSettingsFilmsPageSize       = SettingsAwait + "films_page_size"
	CallSettingsObjectsPageSize      = Settings + "objects_page_size"
	AwaitSettingsObjectsPageSize     = SettingsAwait + "objects_page_size"

	// admin
	Admin              = "admin_select_"
	CallAdminAdmins    = Admin + "admins"
	CallAdminUsers     = Admin + "users"
	CallAdminBroadcast = Admin + "broadcast"
	CallAdminFeedback  = Admin + "feedback"

	// entities
	SelectEntity         = "select_entity_"
	SelectAdmin          = SelectEntity + "admin_"
	SelectUser           = SelectEntity + "user_"
	Entities             = "entities_"
	EntitiesAwait        = Entities + "await_"
	EntitiesPage         = Entities + "page_"
	CallEntitiesBack     = Entities + "back"
	CallEntitiesPagePrev = EntitiesPage + "prev"
	CallEntitiesPageNext = EntitiesPage + "next"
	CallEntitiesPageLast = EntitiesPage + "last"
	CallEntitiesFirst    = EntitiesPage + "first"
	CallEntitiesFind     = Entities + "find"
	AwaitEntitiesFind    = EntitiesAwait + "find"

	// user detail
	UserDetail                     = "user_detail_"
	UserDetailAwait                = UserDetail + "await_"
	UserDetailRole                 = UserDetail + "role_"
	CallUserDetailAgain            = UserDetail + "again"
	CallUserDetailBack             = UserDetail + "back"
	CallUserDetailLogs             = UserDetail + "logs"
	CallUserDetailBan              = UserDetail + "ban"
	AwaitUserDetailReason          = UserDetailAwait + "reason"
	CallUserDetailUnban            = UserDetail + "unban"
	CallUserDetailFeedback         = UserDetail + "feedback"
	CallUserDetailRole             = UserDetail + "role"
	CallUserDetailRoleSelectUser   = UserDetailRole + "user"
	CallUserDetailRoleSelectHelper = UserDetailRole + "helper"
	CallUserDetailRoleSelectAdmin  = UserDetailRole + "admin"
	CallUserDetailRoleSelectSuper  = UserDetailRole + "super"

	// admin detail
	AdminDetail               = "admin_detail_"
	CallAdminDetailBack       = AdminDetail + "back"
	CallAdminDetailAgain      = AdminDetail + "again"
	CallAdminDetailRaiseRole  = AdminDetail + "raise_role"
	CallAdminDetailLowerRole  = AdminDetail + "lower_role"
	CallAdminDetailRemoveRole = AdminDetail + "remove_role"

	// feedbacks
	SelectFeedback         = "select_feedback_"
	Feedbacks              = "feedbacks_"
	FeedbacksPage          = Feedbacks + "page_"
	CallFeedbacksBack      = Feedbacks + "back"
	CallFeedbacksPagePrev  = FeedbacksPage + "prev"
	CallFeedbacksPageNext  = FeedbacksPage + "next"
	CallFeedbacksPageLast  = FeedbacksPage + "last"
	CallFeedbacksPageFirst = FeedbacksPage + "first"

	// feedback detail
	FeedbackDetail           = "feedback_detail_"
	CallFeedbackDetailBack   = FeedbackDetail + "back"
	CallFeedbackDetailDelete = FeedbackDetail + "delete"

	// broadcast
	Broadcast             = "broadcast_"
	BroadcastAwait        = Broadcast + "await_"
	CallBroadcastSend     = Broadcast + "send"
	AwaitBroadcastImage   = BroadcastAwait + "image"
	AwaitBroadcastText    = BroadcastAwait + "text"
	AwaitBroadcastPin     = BroadcastAwait + "pin"
	AwaitBroadcastConfirm = BroadcastAwait + "confirm"

	// profile
	Profile           = "profile_"
	CallProfileUpdate = Profile + "update"
	CallProfileDelete = Profile + "delete"

	// update profile
	UpdateProfile              = "update_profile_"
	UpdateProfileAwait         = UpdateProfile + "await_"
	CallUpdateProfileBack      = UpdateProfile + "back"
	CallUpdateProfileUsername  = UpdateProfile + "username"
	AwaitUpdateProfileUsername = UpdateProfileAwait + "username"
	CallUpdateProfileEmail     = UpdateProfile + "email"
	AwaitUpdateProfileEmail    = UpdateProfileAwait + "email"

	// delete profile
	DeleteProfile             = "delete_profile_"
	DeleteProfileAwait        = DeleteProfile + "await_"
	AwaitDeleteProfileConfirm = DeleteProfileAwait + "confirm"

	// films
	SelectFilm         = "select_film_"
	Films              = "films_"
	FilmsAwait         = Films + "await_"
	FilmsPage          = Films + "page_"
	CallFilmsBack      = Films + "back"
	CallFilmsNew       = Films + "new"
	CallFilmsFind      = Films + "find"
	CallFilmsFilters   = Films + "filters"
	CallFilmsSorting   = Films + "sorting"
	CallFilmsManage    = Films + "manage"
	CallFilmsPageNext  = FilmsPage + "next"
	CallFilmsPagePrev  = FilmsPage + "prev"
	CallFilmsPageLast  = FilmsPage + "last"
	CallFilmsPageFirst = FilmsPage + "first"
	AwaitFilmsTitle    = FilmsAwait + "title"

	// find films
	FindFilms              = "find_films_"
	FindFilmsPage          = FindFilms + "page_"
	CallFindFilmsBack      = FindFilms + "back"
	CallFindFilmsAgain     = FindFilms + "again"
	CallFindFilmsPageNext  = FindFilmsPage + "next"
	CallFindFilmsPagePrev  = FindFilmsPage + "prev"
	CallFindFilmsPageLast  = FindFilmsPage + "last"
	CallFindFilmsPageFirst = FindFilmsPage + "first"

	// find new film
	SelectNewFilm            = "select_new_film_"
	FindNewFilm              = "find_new_film_"
	FindNewFilmPage          = FindNewFilm + "page_"
	CallFindNewFilmBack      = FindNewFilm + "back"
	CallFindNewFilmAgain     = FindNewFilm + "again"
	CallFindNewFilmPageNext  = FindNewFilmPage + "next"
	CallFindNewFilmPagePrev  = FindNewFilmPage + "prev"
	CallFindNewFilmPageLast  = FindNewFilmPage + "last"
	CallFindNewFilmPageFirst = FindNewFilmPage + "first"

	// film filters
	FilmFilters                           = "film_filters_"
	FilmFiltersSelect                     = FilmFilters + "select_"
	FilmFiltersSelectRange                = FilmFiltersSelect + "range_"
	FilmFiltersSelectSwitch               = FilmFiltersSelect + "switch_"
	FilmFiltersAwait                      = FilmFilters + "await_"
	FilmFiltersAwaitRange                 = FilmFiltersAwait + "range_"
	FilmFiltersAwaitSwitch                = FilmFiltersAwait + "switch_"
	CallFilmFiltersBack                   = FilmFilters + "back"
	CallFilmFiltersAllReset               = FilmFilters + "all_reset"
	CallFilmFiltersSelectRangeRating      = FilmFiltersSelectRange + "rating"
	CallFilmFiltersSelectRangeUserRating  = FilmFiltersSelectRange + "user_rating"
	CallFilmFiltersSelectRangeYear        = FilmFiltersSelectRange + "year"
	CallFilmFiltersSelectSwitchIsViewed   = FilmFiltersSelectSwitch + "is_viewed"
	CallFilmFiltersSelectSwitchIsFavorite = FilmFiltersSelectSwitch + "is_favorite"
	CallFilmFiltersSelectSwitchHasURL     = FilmFiltersSelectSwitch + "has_url"

	// film sorting
	FilmSorting                     = "film_sorting_"
	FilmSortingSelect               = FilmSorting + "select_"
	FilmSortingAwait                = FilmSorting + "await_"
	CallFilmSortingBack             = FilmSorting + "back"
	CallFilmSortingAllReset         = FilmSorting + "all_reset"
	AwaitFilmSortingDirection       = FilmSortingAwait + "direction"
	CallFilmSortingSelectTitle      = FilmSortingSelect + "title"
	CallFilmSortingSelectRating     = FilmSortingSelect + "rating"
	CallFilmSortingSelectYear       = FilmSortingSelect + "year"
	CallFilmSortingSelectIsViewed   = FilmSortingSelect + "is_viewed"
	CallFilmSortingSelectIsFavorite = FilmSortingSelect + "is_favorite"
	CallFilmSortingSelectUserRating = FilmSortingSelect + "user_rating"
	CallFilmSortingSelectCreatedAt  = FilmSortingSelect + "created_at"

	// manage film
	ManageFilm                         = "manage_film_"
	CallManageFilmBack                 = ManageFilm + "back"
	CallManageFilmUpdate               = ManageFilm + "update"
	CallManageFilmDelete               = ManageFilm + "delete"
	CallManageFilmRemoveFromCollection = ManageFilm + "remove_from_collection"

	// new film
	NewFilm                         = "new_film_"
	NewFilmAwait                    = NewFilm + "await_"
	CallNewFilmBack                 = NewFilm + "back"
	CallNewFilmManually             = NewFilm + "manually"
	CallNewFilmFromURL              = NewFilm + "from_url"
	CallNewFilmFind                 = NewFilm + "find"
	CallNewFilmChangeKinopoiskToken = NewFilm + "change_kinopoisk_token"
	AwaitNewFilmFind                = NewFilmAwait + "find"
	AwaitNewFilmFromURL             = NewFilmAwait + "from_url"
	AwaitNewFilmTitle               = NewFilmAwait + "title"
	AwaitNewFilmYear                = NewFilmAwait + "year"
	AwaitNewFilmGenre               = NewFilmAwait + "genre"
	AwaitNewFilmDescription         = NewFilmAwait + "description"
	AwaitNewFilmRating              = NewFilmAwait + "rating"
	AwaitNewFilmImage               = NewFilmAwait + "image"
	AwaitNewFilmComment             = NewFilmAwait + "comment"
	AwaitNewFilmViewed              = NewFilmAwait + "viewed"
	AwaitNewFilmFilmURL             = NewFilmAwait + "film_url"
	AwaitNewFilmUserRating          = NewFilmAwait + "user_rating"
	AwaitNewFilmReview              = NewFilmAwait + "review"
	AwaitNewFilmKinopoiskToken      = NewFilmAwait + "kinopoisk_token"

	// delete film
	DeleteFilm             = "delete_film_"
	DeleteFilmAwait        = DeleteFilm + "await_"
	AwaitDeleteFilmConfirm = DeleteFilmAwait + "confirm"

	// update film
	UpdateFilm                 = "update_film_"
	UpdateFilmAwait            = UpdateFilm + "await_"
	CallUpdateFilmBack         = UpdateFilm + "back"
	CallUpdateFilmTitle        = UpdateFilm + "title"
	AwaitUpdateFilmTitle       = UpdateFilmAwait + "title"
	CallUpdateFilmYear         = UpdateFilm + "year"
	AwaitUpdateFilmYear        = UpdateFilmAwait + "year"
	CallUpdateFilmGenre        = UpdateFilm + "genre"
	AwaitUpdateFilmGenre       = UpdateFilmAwait + "genre"
	CallUpdateFilmDescription  = UpdateFilm + "description"
	AwaitUpdateFilmDescription = UpdateFilmAwait + "description"
	CallUpdateFilmRating       = UpdateFilm + "rating"
	AwaitUpdateFilmRating      = UpdateFilmAwait + "rating"
	CallUpdateFilmImage        = UpdateFilm + "image"
	AwaitUpdateFilmImage       = UpdateFilmAwait + "image"
	CallUpdateFilmURL          = UpdateFilm + "url"
	AwaitUpdateFilmURL         = UpdateFilmAwait + "url"
	CallUpdateFilmComment      = UpdateFilm + "comment"
	AwaitUpdateFilmComment     = UpdateFilmAwait + "comment"
	CallUpdateFilmViewed       = UpdateFilm + "viewed"
	AwaitUpdateFilmViewed      = UpdateFilmAwait + "viewed"
	CallUpdateFilmUserRating   = UpdateFilm + "user_rating"
	AwaitUpdateFilmUserRating  = UpdateFilmAwait + "user_rating"
	CallUpdateFilmReview       = UpdateFilm + "review"
	AwaitUpdateFilmReview      = UpdateFilmAwait + "review"

	// film detail
	FilmDetail             = "film_detail_"
	FilmDetailPage         = FilmDetail + "page_"
	CallFilmDetailPageNext = FilmDetailPage + "next"
	CallFilmDetailPagePrev = FilmDetailPage + "prev"
	CallFilmDetailBack     = FilmDetail + "back"
	CallFilmDetailViewed   = FilmDetail + "viewed"
	CallFilmDetailFavorite = FilmDetail + "favorite"

	// viewed film
	ViewedFilm                = "viewed_film_"
	ViewedFilmAwait           = ViewedFilm + "await_"
	AwaitViewedFilmUserRating = ViewedFilmAwait + "user_rating"
	AwaitViewedFilmReview     = ViewedFilmAwait + "review"

	// collections
	SelectCollection         = "select_collection_"
	Collections              = "collections_"
	CollectionsAwait         = Collections + "await_"
	CollectionsPage          = Collections + "page_"
	CallCollectionsNew       = Collections + "new"
	CallCollectionsManage    = Collections + "manage"
	CallCollectionsFavorite  = Collections + "favorite"
	CallCollectionsPageNext  = Collections + "page_next"
	CallCollectionsPagePrev  = Collections + "page_prev"
	CallCollectionsPageLast  = Collections + "page_last"
	CallCollectionsPageFirst = Collections + "page_first"
	CallCollectionsBack      = Collections + "back"
	CallCollectionsFind      = Collections + "find"
	CallCollectionsSorting   = Collections + "sorting"
	AwaitCollectionsName     = CollectionsAwait + "name"

	// collection sorting
	CollectionSorting                     = "collection_sorting_"
	CollectionSortingSelect               = CollectionSorting + "select_"
	CollectionSortingAwait                = CollectionSorting + "await_"
	CallCollectionSortingBack             = CollectionSorting + "back"
	CallCollectionSortingAllReset         = CollectionSorting + "all_reset"
	CallCollectionSortingSelectIsFavorite = CollectionSortingSelect + "is_favorite"
	CallCollectionSortingSelectName       = CollectionSortingSelect + "name"
	CallCollectionSortingSelectCreatedAt  = CollectionSortingSelect + "created_at"
	CallCollectionSortingSelectTotalFilms = CollectionSortingSelect + "total_films"
	AwaitCollectionSortingDirection       = CollectionSortingAwait + "direction"

	// find collections
	FindCollections              = "find_collections_"
	FindCollectionsPage          = FindCollections + "page_"
	CallFindCollectionsBack      = FindCollections + "back"
	CallFindCollectionsPageNext  = FindCollectionsPage + "next"
	CallFindCollectionsPagePrev  = FindCollectionsPage + "prev"
	CallFindCollectionsPageLast  = FindCollectionsPage + "last"
	CallFindCollectionsPageFirst = FindCollectionsPage + "first"
	CallFindCollectionsAgain     = FindCollections + "again"

	// manage collection
	ManageCollection           = "manage_collection_"
	CallManageCollectionBack   = ManageCollection + "back"
	CallManageCollectionUpdate = ManageCollection + "update"
	CallManageCollectionDelete = ManageCollection + "delete"

	// new collection
	NewCollection                 = "new_collection_"
	NewCollectionAwait            = NewCollection + "await_"
	AwaitNewCollectionName        = NewCollectionAwait + "name"
	AwaitNewCollectionDescription = NewCollectionAwait + "description"

	// delete collection
	DeleteCollection             = "delete_collection_"
	DeleteCollectionAwait        = DeleteCollection + "await_"
	AwaitDeleteCollectionConfirm = DeleteCollectionAwait + "confirm"

	// update colllection
	UpdateCollection                 = "update_collection_"
	UpdateCollectionAwait            = UpdateCollection + "await_"
	CallUpdateCollectionBack         = UpdateCollection + "back"
	CallUpdateCollectionName         = UpdateCollection + "name"
	AwaitUpdateCollectionName        = UpdateCollectionAwait + "name"
	CallUpdateCollectionDescription  = UpdateCollection + "description"
	AwaitUpdateCollectionDescription = UpdateCollectionAwait + "description"

	// collection films
	CollectionFilmsFrom               = "collection_films_from_"
	CallCollectionFilmsFromFilm       = CollectionFilmsFrom + "film"
	CallCollectionFilmsFromCollection = CollectionFilmsFrom + "collection"

	// film to collection option
	FilmToCollectionOption             = "film_to_collection_option_"
	CallFilmToCollectionOptionBack     = FilmToCollectionOption + "back"
	CallFilmToCollectionOptionNew      = FilmToCollectionOption + "new"
	CallFilmToCollectionOptionExisting = FilmToCollectionOption + "existing"

	// add collection to film
	SelectCFCollection               = "select_cf_collection_"
	AddCollectionToFilm              = "add_collection_to_film_"
	AddCollectionToFilmPage          = AddCollectionToFilm + "page_"
	AddCollectionToFilmAwait         = AddCollectionToFilm + "await_"
	CallAddCollectionToFilmBack      = AddCollectionToFilm + "back"
	CallAddCollectionToFilmPagePrev  = AddCollectionToFilmPage + "prev"
	CallAddCollectionToFilmPageNext  = AddCollectionToFilmPage + "next"
	CallAddCollectionToFilmPageLast  = AddCollectionToFilmPage + "last"
	CallAddCollectionToFilmPageFirst = AddCollectionToFilmPage + "first"
	CallAddCollectionToFilmFind      = AddCollectionToFilm + "find"
	CallAddCollectionToFilmAgain     = AddCollectionToFilm + "again"
	CallAddCollectionToFilmReset     = AddCollectionToFilm + "reset"
	AwaitAddCollectionToFilmName     = AddCollectionToFilmAwait + "name"

	// add film to collection
	SelectCFFilm                     = "select_cf_film_"
	AddFilmToCollection              = "add_film_to_collection_"
	AddFilmToCollectionPage          = AddFilmToCollection + "page_"
	AddFilmToCollectionAwait         = AddFilmToCollectionPage + "await_"
	CallAddFilmToCollectionBack      = AddFilmToCollection + "back"
	CallAddFilmToCollectionPagePrev  = AddFilmToCollectionPage + "prev"
	CallAddFilmToCollectionPageNext  = AddFilmToCollectionPage + "next"
	CallAddFilmToCollectionPageLast  = AddFilmToCollectionPage + "last"
	CallAddFilmToCollectionPageFirst = AddFilmToCollectionPage + "first"
	CallAddFilmToCollectionFind      = AddFilmToCollection + "find"
	CallAddFilmToCollectionAgain     = AddFilmToCollection + "again"
	CallAddFilmToCollectionReset     = AddFilmToCollection + "reset"
	AwaitAddFilmToCollectionTitle    = AddFilmToCollectionAwait + "title"
)
