package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

func HandleFiltersFilmsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFiltersFilmsMessage(session), keyboards.BuildFilmsFilterKeyboard(session))
}

func HandleFiltersFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackFiltersFilmsBack:
		HandleFilmsCommand(app, session)
	case states.CallbackFiltersFilmsAllReset:
		handleFiltersFilmsAllReset(app, session)
	default:
		if strings.HasPrefix(callback, states.PrefixFiltersFilmsSelect) {
			handleFiltersFilmsSelect(app, session, callback)
		}
	}
}

func HandleFiltersFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		resetFilmsStateAndHandleFiltersFilms(app, session)
		return
	}

	switch {
	case strings.HasPrefix(session.State, states.PrefixFiltersFilmsAwaitingRange):
		parseFiltersFilmsRange(app, session, strings.TrimPrefix(session.State, states.PrefixFiltersFilmsAwaitingRange))
	case strings.HasPrefix(session.State, states.PrefixFiltersFilmsAwaitingSwitch):
		parseFiltersFilmsSwitch(app, session, strings.TrimPrefix(session.State, states.PrefixFiltersFilmsAwaitingSwitch))
	}
}

func handleFiltersFilmsSelect(app models.App, session *models.Session, callback string) {
	switch {
	case strings.HasPrefix(callback, states.PrefixFiltersFilmsSelectRange):
		handleFiltersFilmsRange(app, session, strings.TrimPrefix(callback, states.PrefixFiltersFilmsSelectRange))
	case strings.HasPrefix(callback, states.PrefixFiltersFilmsSelectSwitch):
		handleFiltersFilmsSwitch(app, session, strings.TrimPrefix(callback, states.PrefixFiltersFilmsSelectSwitch))
	}
}

func handleFiltersFilmsAllReset(app models.App, session *models.Session) {
	session.GetFilmsFiltersByContext().ResetFilters()
	app.SendMessage(messages.BuildFilterResetSuccessSimpleMessage(session), nil)
	resetFilmsStateAndHandleFiltersFilms(app, session)
}

func handleFiltersFilmsSwitch(app models.App, session *models.Session, filterType string) {
	app.SendMessage(messages.BuildFilterSwitchMessage(session, filterType), keyboards.BuildFiltersFilmsSwitchKeyboard(session, filterType))
	session.SetState(states.PrefixFiltersFilmsAwaitingSwitch + filterType)
}

func parseFiltersFilmsSwitch(app models.App, session *models.Session, filterType string) {
	if utils.IsReset(app.Update) {
		handleFiltersFilmsReset(app, session, filterType)
		return
	}

	session.GetFilmsFiltersByContext().ApplySwitchValue(filterType, utils.IsAgree(app.Update))
	handleFiltersFilmsApplied(app, session, filterType, "🔀")
}

func handleFiltersFilmsRange(app models.App, session *models.Session, filterType string) {
	app.SendMessage(messages.BuildFilterRangeMessage(session, filterType), keyboards.BuildFiltersFilmsRangeKeyboard(session, filterType))
	session.SetState(states.PrefixFiltersFilmsAwaitingRange + filterType)
}

func parseFiltersFilmsRange(app models.App, session *models.Session, filterType string) {
	if utils.IsReset(app.Update) {
		handleFiltersFilmsReset(app, session, filterType)
		return
	}

	config := getFilterRangeConfig(filterType)
	input, err := utils.ValidateFiltersRange(utils.ParseMessageString(app.Update), config)
	if err != nil {
		handleFiltersInvalidRangeInput(app, session, filterType, config)
		return
	}

	session.GetFilmsFiltersByContext().ApplyRangeValue(filterType, input)
	handleFiltersFilmsApplied(app, session, filterType, "↕️")
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

func handleFiltersFilmsReset(app models.App, session *models.Session, filterType string) {
	session.GetFilmsFiltersByContext().ResetFilter(filterType)
	app.SendMessage(messages.BuildFilterResetSuccessMessage(session, filterType), nil)
	resetFilmsStateAndHandleFiltersFilms(app, session)
}

func handleFiltersFilmsApplied(app models.App, session *models.Session, filterType, emoji string) {
	app.SendMessage(messages.BuildFilterAppliedMessage(session, filterType, emoji), nil)
	resetFilmsStateAndHandleFiltersFilms(app, session)
}

func handleFiltersInvalidRangeInput(app models.App, session *models.Session, filterType string, config utils.FilterRangeConfig) {
	app.SendMessage(messages.BuildInvalidFilterRangeInputMessage(session, config), nil)
	handleFiltersFilmsRange(app, session, filterType)
}

func resetFilmsStateAndHandleFiltersFilms(app models.App, session *models.Session) {
	session.FilmsState.CurrentPage = 1
	session.ClearState()
	HandleFiltersFilmsCommand(app, session)
}
