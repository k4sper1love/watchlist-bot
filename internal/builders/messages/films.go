package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildFilmsMessage(session *models.Session, metadata *filters.Metadata) string {
	var header string
	switch session.Context {
	case states.ContextFilm:
		header = ""
	case states.ContextCollection:
		header = BuildCollectionHeader(session)
	default:
		return translator.Translate(session.Lang, "unknownContext", nil, nil)
	}

	return header + buildFilmsList(session, metadata, false, true)
}

func BuildFindFilmsMessage(session *models.Session, metadata *filters.Metadata) string {
	return buildFilmsList(session, metadata, true, true)
}

func BuildFindNewFilmMessage(session *models.Session, metadata *filters.Metadata) string {
	return buildFilmsList(session, metadata, true, false)
}

func buildFilmsList(session *models.Session, metadata *filters.Metadata, isFind bool, needViewed bool) string {
	films := session.FilmsState.Films

	if metadata.TotalRecords == 0 {
		msg := "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
		return msg
	}

	totalFilmsMsgKey := "totalFilms"
	if isFind {
		totalFilmsMsgKey = "totalFindFilms"
	}

	totalFilmsMsg := translator.Translate(session.Lang, totalFilmsMsgKey, nil, nil)
	msg := fmt.Sprintf("🎥 <b>%s:</b> %d\n\n", totalFilmsMsg, metadata.TotalRecords)

	for i, film := range films {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)
		numberEmoji := utils.NumberToEmoji(itemID)
		msg += fmt.Sprintf("%s ", numberEmoji)

		if film.IsFavorite {
			msg += "⭐ "
		}

		msg += fmt.Sprintf(" <i>ID: %d</i>", film.ID)

		msg += "\n" + BuildFilmGeneralMessage(session, &film, needViewed)
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": metadata.CurrentPage,
		"LastPage":    metadata.LastPage,
	}, nil)
	msg += fmt.Sprintf("<b>📄 %s</b>", pageMsg)

	return msg
}

func BuildCollectionHeader(session *models.Session) string {
	collection := session.CollectionDetailState.Collection

	msg := fmt.Sprintf("<b>%s</b>", collection.Name)

	if collection.IsFavorite {
		msg += " ⭐"
	}

	if collection.Description != "" {
		msg += fmt.Sprintf("\n<i>%s</i>", collection.Description)
	}

	msg += "\n\n"

	return msg
}

func BuildFilterRangeMessage(session *models.Session, filterType string) string {
	filter := session.GetFilmsFiltersByContext()

	part1 := translator.Translate(session.Lang, "filterInstructionRange", nil, nil)
	part2 := translator.Translate(session.Lang, "filterInstructionPartialRange", nil, nil)

	msg := fmt.Sprintf("↕️ %s\n\n<i>%s</i>", part1, part2)

	if filter.IsFilterEnabled(filterType) {
		currentValueMsg := translator.Translate(session.Lang, "currentValue", nil, nil)
		value := filter.ValueToString(filterType)
		msg += fmt.Sprintf("\n\n<b>%s</b>: %s", currentValueMsg, value)
	}

	return msg
}

func BuildFilterSwitchMessage(session *models.Session, filterType string) string {
	filter := session.GetFilmsFiltersByContext()

	filterMsg := translator.Translate(session.Lang, filterType, nil, nil)
	msg := "🔀 " + translator.Translate(session.Lang, "filterInstructionSwitch", map[string]interface{}{
		"Filter": filterMsg,
	}, nil)

	if filter.IsFilterEnabled(filterType) {
		currentValueMsg := translator.Translate(session.Lang, "currentValue", nil, nil)
		value := translator.Translate(session.Lang, filter.ValueToString(filterType), nil, nil)
		msg += fmt.Sprintf("\n\n<b>%s</b>: %s", currentValueMsg, value)
	}

	return msg
}

func BuildInvalidFilterRangeInputMessage(session *models.Session, config models.FilterRangeConfig) string {
	exampleValue := translator.Translate(session.Lang, "exampleValue", nil, nil)
	exampleRange := translator.Translate(session.Lang, "exampleRange", nil, nil)
	examplePartialRange := translator.Translate(session.Lang, "examplePartialRange", nil, nil)
	rangeLimits := translator.Translate(session.Lang, "rangeLimits", map[string]interface{}{
		"Min": fmt.Sprintf("%.f", config.MinValue),
		"Max": fmt.Sprintf("%.f", config.MaxValue),
	}, nil)

	msg := "❌ " + translator.Translate(session.Lang, "invalidInput", nil, nil)
	msg += "\n\n<b>" + translator.Translate(session.Lang, "requestRangeInFormat", nil, nil) + "</b>"
	msg += fmt.Sprintf("\n- %s: <code>%s</code>", exampleValue, "5.5")
	msg += fmt.Sprintf("\n- %s: <code>%s</code>", exampleRange, "1990-2023")
	msg += fmt.Sprintf("\n- %s: <code>%s</code> или <code>%s</code>", examplePartialRange, "5-", "-10")
	msg += fmt.Sprintf("\n\n⚠️ <i>%s</i>", rangeLimits)

	return msg
}
