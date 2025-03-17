package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

// HandleFilmsCommand handles the command for listing films.
// Retrieves paginated films based on the current session context (user or collection)
// and sends a message with their details and navigation buttons.
func HandleFilmsCommand(app models.App, session *models.Session) {
	// Clears the title used for finding films in other contexts to ensure a fresh state.
	session.FilmsState.Clear()

	if metadata, err := getFilms(app, session); err != nil {
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, ""))
	} else {
		app.SendMessage(messages.Films(session, metadata), keyboards.Films(session, metadata.CurrentPage, metadata.LastPage))
	}
}

// HandleFilmsButtons handles button interactions related to the films list.
// Supports actions like going back, creating new films, managing existing ones, searching, filtering, sorting, and pagination.
func HandleFilmsButtons(app models.App, session *models.Session, back func(models.App, *models.Session)) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFilmsBack:
		if session.Context == states.CtxCollection {
			session.CollectionsState.CurrentPage = 1
		}
		back(app, session)

	case states.CallFilmsNew:
		HandleNewFilmCommand(app, session)

	case states.CallFilmsManage:
		HandleManageFilmCommand(app, session)

	case states.CallFilmsFind:
		requestFindFilmsTitle(app, session)

	case states.CallFilmsFilters:
		HandleFilmFiltersCommand(app, session)

	case states.CallFilmsSorting:
		HandleSortingFilmsCommand(app, session)

	default:
		if strings.HasPrefix(callback, states.FilmsPage) {
			handleFilmsPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectFilm) {
			handleFilmSelect(app, session, callback)
		}
	}
}

// HandleFilmsProcess processes workflows related to the films list.
// Handles states like awaiting a film title input for search.
func HandleFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitFilmsTitle:
		parser.ParseFindFilmsTitle(app, session, HandleFindFilmsCommand)
	}
}

// handleFilmsPagination processes pagination actions for the films list.
// Updates the current page in the session and reloads the films list.
func handleFilmsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallFilmsPageNext:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallFilmsPagePrev:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallFilmsPageLast:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallFilmsPageFirst:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1
	}

	HandleFilmsCommand(app, session)
}

// handleFilmSelect processes the selection of a film from the list.
// Parses the film index and navigates to the detailed view of the selected film.
func handleFilmSelect(app models.App, session *models.Session, callback string) {
	if index, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectFilm)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallFilmsBack))
	} else {
		session.FilmDetailState.Index = index
		HandleFilmDetailCommand(app, session)
	}
}

// requestFindFilmsTitle prompts the user to enter the title of a film to search for.
func requestFindFilmsTitle(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.AwaitFilmsTitle)
}

// updateFilmList updates the film list by navigating to the next or previous page.
func updateFilmList(app models.App, session *models.Session, next bool) error {
	if next {
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			return fmt.Errorf(messages.LastPageAlert(session))
		}
		session.FilmsState.CurrentPage++

	} else {
		if session.FilmsState.CurrentPage <= 1 {
			return fmt.Errorf(messages.FirstPageAlert(session))
		}
		session.FilmsState.CurrentPage--
	}

	_, err := getFilms(app, session)
	return err
}

// getFilms retrieves a paginated list of films based on the session context (user or collection).
// Updates the session with the retrieved films and their metadata.
func getFilms(app models.App, session *models.Session) (*filters.Metadata, error) {
	var (
		films    []apiModels.Film
		metadata *filters.Metadata
		err      error
	)

	switch session.Context {
	case states.CtxFilm:
		films, metadata, err = fetchFilmsFromUser(app, session)
	case states.CtxCollection:
		films, metadata, err = fetchFilmsFromCollection(app, session)
	default:
		return nil, fmt.Errorf("unsupported session context: %v", session.Context)
	}

	if err != nil {
		return nil, err
	}

	updateSessionWithFilms(session, films, metadata)
	return metadata, nil
}

// fetchFilmsFromUser retrieves films associated with the current user using the Watchlist service.
func fetchFilmsFromUser(app models.App, session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	filmsResponse, err := watchlist.GetFilms(app, session)
	if err != nil {
		return nil, nil, err
	}

	return filmsResponse.Films, &filmsResponse.Metadata, nil
}

// fetchFilmsFromCollection retrieves films associated with the current collection using the Watchlist service.
func fetchFilmsFromCollection(app models.App, session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	collectionResponse, err := watchlist.GetCollectionFilms(app, session)
	if err != nil {
		return nil, nil, err
	}

	session.CollectionDetailState.Collection = collectionResponse.CollectionFilms.Collection
	return collectionResponse.CollectionFilms.Films, &collectionResponse.Metadata, nil
}

// updateSessionWithFilms updates the session with the retrieved films and their metadata.
func updateSessionWithFilms(session *models.Session, films []apiModels.Film, metadata *filters.Metadata) {
	session.FilmsState.Films = films
	session.FilmsState.LastPage = metadata.LastPage
	session.FilmsState.TotalRecords = metadata.TotalRecords
	session.FilmsState.PageSize = metadata.PageSize
}
