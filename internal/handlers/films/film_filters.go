package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

func HandleFilmFiltersCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceFilter(session), keyboards.FilmFilters(session))
}

func HandleFilmFiltersButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFilmFiltersBack:
		HandleFilmsCommand(app, session)

	case states.CallFilmFiltersAllReset:
		handleFilmFiltersAllReset(app, session)

	default:
		if strings.HasPrefix(callback, states.FilmFiltersSelect) {
			handleFilmFiltersSelect(app, session, callback)
		}
	}
}

func HandleFilmFiltersProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		resetFilmsStateAndHandleFilmFilters(app, session)
		return
	}

	switch {
	case strings.HasPrefix(session.State, states.FilmFiltersAwaitRange):
		parseFilmFiltersRange(app, session, strings.TrimPrefix(session.State, states.FilmFiltersAwaitRange))
	case strings.HasPrefix(session.State, states.FilmFiltersAwaitSwitch):
		parseFilmFiltersSwitch(app, session, strings.TrimPrefix(session.State, states.FilmFiltersAwaitSwitch))
	}
}

func handleFilmFiltersSelect(app models.App, session *models.Session, callback string) {
	switch {
	case strings.HasPrefix(callback, states.FilmFiltersSelectRange):
		handleFilmFiltersRange(app, session, strings.TrimPrefix(callback, states.FilmFiltersSelectRange))
	case strings.HasPrefix(callback, states.FilmFiltersSelectSwitch):
		handleFilmFiltersSwitch(app, session, strings.TrimPrefix(callback, states.FilmFiltersSelectSwitch))
	}
}

func handleFilmFiltersAllReset(app models.App, session *models.Session) {
	session.GetFilmFiltersByCtx().ResetAll()
	app.SendMessage(messages.ResetFiltersSuccess(session), nil)
	resetFilmsStateAndHandleFilmFilters(app, session)
}

func handleFilmFiltersSwitch(app models.App, session *models.Session, filterType string) {
	app.SendMessage(messages.FilterSwitch(session, filterType), keyboards.FilmFilterSwitch(session, filterType))
	session.SetState(states.FilmFiltersAwaitSwitch + filterType)
}

func parseFilmFiltersSwitch(app models.App, session *models.Session, filterType string) {
	if utils.IsReset(app.Update) {
		handleFilmFiltersReset(app, session, filterType)
		return
	}

	session.GetFilmFiltersByCtx().ApplySwitch(filterType, utils.IsAgree(app.Update))
	handleFilmFiltersApplied(app, session, filterType, "🔀")
}

func handleFilmFiltersRange(app models.App, session *models.Session, filterType string) {
	app.SendMessage(messages.FilterRange(session, filterType), keyboards.FilmFilterRange(session, filterType))
	session.SetState(states.FilmFiltersAwaitRange + filterType)
}

func parseFilmFiltersRange(app models.App, session *models.Session, filterType string) {
	if utils.IsReset(app.Update) {
		handleFilmFiltersReset(app, session, filterType)
		return
	}

	config := getFilterRangeConfig(filterType)
	input, err := utils.ValidateFiltersRange(utils.ParseMessageString(app.Update), config)
	if err != nil {
		handleFiltersInvalidRangeInput(app, session, filterType, config)
		return
	}

	session.GetFilmFiltersByCtx().ApplyRange(filterType, input)
	handleFilmFiltersApplied(app, session, filterType, "↕️")
}

func getFilterRangeConfig(filterType string) utils.FilterRangeConfig {
	switch filterType {
	case "rating", "user_rating":
		return utils.FilterRangeConfig{MinValue: 0, MaxValue: 10}
	case "year":
		return utils.FilterRangeConfig{MinValue: 1888, MaxValue: 2100}
	default:
		return utils.FilterRangeConfig{}
	}
}

func handleFilmFiltersReset(app models.App, session *models.Session, filterType string) {
	session.GetFilmFiltersByCtx().Reset(filterType)
	app.SendMessage(messages.ResetFilterSuccess(session, filterType), nil)
	resetFilmsStateAndHandleFilmFilters(app, session)
}

func handleFilmFiltersApplied(app models.App, session *models.Session, filterType, emoji string) {
	app.SendMessage(messages.FilterApplied(session, filterType, emoji), nil)
	resetFilmsStateAndHandleFilmFilters(app, session)
}

func handleFiltersInvalidRangeInput(app models.App, session *models.Session, filterType string, config utils.FilterRangeConfig) {
	app.SendMessage(messages.InvalidFilterRange(session, config), nil)
	handleFilmFiltersRange(app, session, filterType)
}

func resetFilmsStateAndHandleFilmFilters(app models.App, session *models.Session) {
	session.FilmsState.CurrentPage = 1
	session.ClearState()
	HandleFilmFiltersCommand(app, session)
}
