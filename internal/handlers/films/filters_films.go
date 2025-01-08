package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleFiltersFilmsCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "choiceFilter", nil, nil)

	keyboard := keyboards.BuildFilmsFilterKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleFiltersFilmsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackFiltersFilmsSelectBack:
		HandleFilmsCommand(app, session)

	case states.CallbackFiltersFilmsSelectAllReset:
		handleFiltersFilmsAllReset(app, session)

	case states.CallbackFiltersFilmsSelectMinRating:
		handleFiltersFilmsRating(app, session, "minRating")

	case states.CallbackFiltersFilmsSelectMaxRating:
		handleFiltersFilmsRating(app, session, "maxRating")
	}
}

func HandleFiltersFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleFiltersFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFiltersFilmsAwaitingMinRating:
		parseFiltersFilmsRating(app, session, "minRating")
	case states.ProcessFiltersFilmsAwaitingMaxRating:
		parseFiltersFilmsRating(app, session, "maxRating")
	}
}

func handleFiltersFilmsAllReset(app models.App, session *models.Session) {
	session.GetFilmsFiltersByContext().ResetFilters()

	msg := translator.Translate(session.Lang, "filterResetSuccess", nil, nil)

	app.SendMessage(msg, nil)

	HandleFiltersFilmsCommand(app, session)
}

func handleFiltersFilmsRating(app models.App, session *models.Session, filterType string) {
	msg := messages.BuildRatingFilterMessage(session, filterType)

	keyboard := keyboards.NewKeyboard().AddResetFilter(session, filterType).AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(getAwaitingState(filterType))
}

func parseFiltersFilmsRating(app models.App, session *models.Session, filterType string) {
	if utils.IsReset(app.Upd) {
		resetFilter(session, filterType)
		handleFiltersFilmsReset(app, session)
		return
	}

	rating := utils.ParseMessageFloat(app.Upd)

	if rating <= 0 || rating >= 10 {
		msg := translator.Translate(session.Lang, "badRating", nil, nil)
		app.SendMessage(msg, nil)
		handleFiltersFilmsRating(app, session, filterType)
		return
	}

	if ok := validateRatingBounds(session, filterType, rating); !ok {
		msg := messages.BuildValidateFilterMessage(session, filterType)
		app.SendMessage(msg, nil)
		handleFiltersFilmsRating(app, session, filterType)
		return
	}

	applyRatingFilter(session, filterType, rating)

	msg := translator.Translate(session.Lang, "filterApplied", nil, nil)
	app.SendMessage(msg, nil)

	session.ClearState()

	HandleFiltersFilmsCommand(app, session)
}

func handleFiltersFilmsReset(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "filterResetSuccess", nil, 1)
	app.SendMessage(msg, nil)

	session.ClearState()

	HandleFiltersFilmsCommand(app, session)
}

func resetFilter(session *models.Session, filterType string) {
	filter := session.GetFilmsFiltersByContext()

	switch filterType {
	case "minRating":
		filter.MinRating = 0
	case "maxRating":
		filter.MaxRating = 0
	}
}

func validateRatingBounds(session *models.Session, filterType string, rating float64) bool {
	switch filterType {
	case "minRating":
		maxRating := session.GetFilmsFiltersByContext().MaxRating
		return maxRating == 0 || rating <= maxRating
	case "maxRating":
		minRating := session.GetFilmsFiltersByContext().MinRating
		return minRating == 0 || rating >= minRating
	}
	return false
}

func applyRatingFilter(session *models.Session, filterType string, rating float64) {
	filter := session.GetFilmsFiltersByContext()

	switch filterType {
	case "minRating":
		filter.MinRating = rating
	case "maxRating":
		filter.MaxRating = rating
	}
}

func getAwaitingState(filterType string) string {
	switch filterType {
	case "minRating":
		return states.ProcessFiltersFilmsAwaitingMinRating
	case "maxRating":
		return states.ProcessFiltersFilmsAwaitingMaxRating
	}
	return ""
}
