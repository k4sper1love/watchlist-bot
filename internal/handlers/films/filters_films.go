package films

import (
	"errors"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"regexp"
	"strconv"
	"strings"
)

func HandleFiltersFilmsCommand(app models.App, session *models.Session) {
	choiceMsg := translator.Translate(session.Lang, "choiceFilter", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>", choiceMsg)

	keyboard := keyboards.BuildFilmsFilterKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleFiltersFilmsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackFiltersFilmsSelectBack:
		HandleFilmsCommand(app, session)

	case states.CallbackFiltersFilmsSelectAllReset:
		handleFiltersFilmsAllReset(app, session)

	case states.CallbackFiltersFilmsSelectRating:
		handleFiltersFilmsRange(app, session, "rating")

	case states.CallbackFiltersFilmsSelectUserRating:
		handleFiltersFilmsRange(app, session, "userRating")

	case states.CallbackFiltersFilmsSelectYear:
		handleFiltersFilmsRange(app, session, "year")

	case states.CallbackFiltersFilmsSelectIsViewed:
		handleFiltersFilmsSwitch(app, session, "isViewed")

	case states.CallbackFiltersFilmsSelectIsFavorite:
		handleFiltersFilmsSwitch(app, session, "isFavorite")

	case states.CallbackFiltersFilmsSelectHasURL:
		handleFiltersFilmsSwitch(app, session, "hasURL")
	}
}

func HandleFiltersFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleFiltersFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFiltersFilmsAwaitingRating:
		parseFiltersFilmsRange(app, session, "rating")
	case states.ProcessFiltersFilmsAwaitingUserRating:
		parseFiltersFilmsRange(app, session, "userRating")
	case states.ProcessFiltersFilmsAwaitingYear:
		parseFiltersFilmsRange(app, session, "year")
	case states.ProcessFiltersFilmsAwaitingIsViewed:
		parseFiltersFilmsSwitch(app, session, "isViewed")
	case states.ProcessFiltersFilmsAwaitingIsFavorite:
		parseFiltersFilmsSwitch(app, session, "isFavorite")
	case states.ProcessFiltersFilmsAwaitingHasURL:
		parseFiltersFilmsSwitch(app, session, "hasURL")
	}
}

func handleFiltersFilmsAllReset(app models.App, session *models.Session) {
	session.GetFilmsFiltersByContext().ResetFilters()

	msg := "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", nil, nil)

	app.SendMessage(msg, nil)

	session.FilmsState.CurrentPage = 1
	HandleFiltersFilmsCommand(app, session)
}

func handleFiltersFilmsSwitch(app models.App, session *models.Session, filterType string) {
	msg := messages.BuildFilterSwitchMessage(session, filterType)

	keyboard := keyboards.NewKeyboard().AddSurvey().AddResetFilmsFilter(session, filterType).AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(getAwaitingState(filterType))
}

func parseFiltersFilmsSwitch(app models.App, session *models.Session, filterType string) {
	filter := session.GetFilmsFiltersByContext()

	if utils.IsReset(app.Upd) {
		filter.ResetFilter(filterType)
		handleFiltersFilmsReset(app, session, filterType)
		return
	}

	value := utils.IsAgree(app.Upd)

	filter.ApplySwitchValue(filterType, value)

	handleFiltersFilmsApplied(app, session, filterType, "🔀 ")
}

func handleFiltersFilmsRange(app models.App, session *models.Session, filterType string) {
	msg := messages.BuildFilterRangeMessage(session, filterType)

	keyboard := keyboards.NewKeyboard().AddResetFilmsFilter(session, filterType).AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(getAwaitingState(filterType))
}

func parseFiltersFilmsRange(app models.App, session *models.Session, filterType string) {
	filter := session.GetFilmsFiltersByContext()

	if utils.IsReset(app.Upd) {
		filter.ResetFilter(filterType)
		handleFiltersFilmsReset(app, session, filterType)
		return
	}

	input := utils.ParseMessageString(app.Upd)

	var err error
	switch filterType {
	case "rating", "userRating":
		config := models.FilterRangeConfig{MinValue: 0, MaxValue: 10}
		input, err = validateFiltersRange(input, config)
		if err != nil {
			handleFiltersInvalidRangeInput(app, session, filterType, config)
			return
		}

	case "year":
		config := models.FilterRangeConfig{MinValue: 1888, MaxValue: 2100}
		input, err = validateFiltersRange(input, config)
		if err != nil {
			handleFiltersInvalidRangeInput(app, session, filterType, config)
			return
		}

	default:
		msg := "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackFilmsFilters).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		session.ClearState()
		return
	}

	filter.ApplyRangeValue(filterType, input)

	handleFiltersFilmsApplied(app, session, filterType, "↕️ ")
}

func handleFiltersFilmsReset(app models.App, session *models.Session, filterType string) {
	filterMsg := translator.Translate(session.Lang, filterType, nil, nil)
	msg := "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", map[string]interface{}{
		"Filter": filterMsg,
	}, 1)

	app.SendMessage(msg, nil)

	session.ClearState()

	session.FilmsState.CurrentPage = 1
	HandleFiltersFilmsCommand(app, session)
}

func handleFiltersFilmsApplied(app models.App, session *models.Session, filterType, emoji string) {
	filterMsg := translator.Translate(session.Lang, filterType, nil, nil)
	msg := emoji + translator.Translate(session.Lang, "filterApplied", map[string]interface{}{
		"Filter": filterMsg,
	}, nil)

	app.SendMessage(msg, nil)

	session.ClearState()

	session.FilmsState.CurrentPage = 1
	HandleFiltersFilmsCommand(app, session)
}

func handleFiltersInvalidRangeInput(app models.App, session *models.Session, filterType string, config models.FilterRangeConfig) {
	msg := messages.BuildInvalidFilterRangeInputMessage(session, config)
	app.SendMessage(msg, nil)

	session.ClearState()

	handleFiltersFilmsRange(app, session, filterType)
}

func getAwaitingState(filterType string) string {
	switch filterType {
	case "rating":
		return states.ProcessFiltersFilmsAwaitingRating
	case "userRating":
		return states.ProcessFiltersFilmsAwaitingUserRating
	case "year":
		return states.ProcessFiltersFilmsAwaitingYear
	case "isViewed":
		return states.ProcessFiltersFilmsAwaitingIsViewed
	case "isFavorite":
		return states.ProcessFiltersFilmsAwaitingIsFavorite
	case "hasURL":
		return states.ProcessFiltersFilmsAwaitingHasURL
	}
	return ""
}

func validateFiltersRange(input string, config models.FilterRangeConfig) (string, error) {
	singleValuePattern := `^\d+(\.\d+)?$`                        // Например: 5, 5.5
	rangeValuePattern := `^\d+(\.\d+)?-\d+(\.\d+)?$`             // Например: 5.5-7.3 или 1990-2023
	incompleteRangePattern := `^(-?\d+(\.\d+)?-|-?\d+(\.\d+)?)$` // Например: 5.5-, -7.3

	input = strings.TrimSpace(input)

	matchSingle, _ := regexp.MatchString(singleValuePattern, input)
	if matchSingle {
		value, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return "", fmt.Errorf("invalid number: %s", input)
		}
		if value < config.MinValue || value > config.MaxValue {
			return "", fmt.Errorf("value out of range: %s (must be between %.2f and %.2f)", input, config.MinValue, config.MaxValue)
		}
		return input, nil
	}

	matchRange, _ := regexp.MatchString(rangeValuePattern, input)
	if matchRange {
		parts := strings.Split(input, "-")
		if len(parts) == 2 {
			start, err1 := strconv.ParseFloat(parts[0], 64)
			end, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 != nil || err2 != nil {
				return "", errors.New("invalid range format")
			}
			if start >= end {
				return "", fmt.Errorf("invalid range: start (%s) must be less than end (%s)", parts[0], parts[1])
			}
			if start < config.MinValue || end > config.MaxValue {
				return "", fmt.Errorf("range out of bounds: %s (must be between %.2f and %.2f)", input, config.MinValue, config.MaxValue)
			}
			return input, nil
		}
	}

	matchIncompleteRange, _ := regexp.MatchString(incompleteRangePattern, input)
	if matchIncompleteRange {
		parts := strings.Split(input, "-")
		if len(parts) == 2 {
			if parts[0] != "" && parts[1] == "" {
				start, err := strconv.ParseFloat(parts[0], 64)
				if err != nil {
					return "", fmt.Errorf("invalid start value: %s", parts[0])
				}
				if start < config.MinValue || start > config.MaxValue {
					return "", fmt.Errorf("value out of range: %s (must be between %.2f and %.2f)", parts[0], config.MinValue, config.MaxValue)
				}

				input = fmt.Sprintf("%s-%.f", parts[0], config.MaxValue)
				return input, nil
			}

			if parts[0] == "" && parts[1] != "" {
				end, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					return "", fmt.Errorf("invalid end value: %s", parts[1])
				}
				if end < config.MinValue || end > config.MaxValue {
					return "", fmt.Errorf("value out of range: %s (must be between %.2f and %.2f)", parts[1], config.MinValue, config.MaxValue)
				}

				input = fmt.Sprintf("%.f-%s", config.MinValue, parts[1])
				return input, nil
			}
		}
	}

	return "", fmt.Errorf("invalid filter format: %s", input)
}
