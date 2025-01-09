package messages

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildFilmDetailMessage(session *models.Session, film *apiModels.Film) string {
	msg := ""

	if session.Context == states.ContextCollection {
		msg += fmt.Sprintf(" <code>%s</code>", session.CollectionDetailState.Collection.Name)
	}

	msg += "\n"

	if film.Title != "" {
		part := translator.Translate(session.Lang, "title", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>: %s\n", part, film.Title)
	}

	part := translator.Translate(session.Lang, "viewed", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>: %s\n", part, utils.BoolToEmoji(film.IsViewed))

	if film.Genre != "" {
		part = translator.Translate(session.Lang, "genre", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>: %s\n", part, film.Genre)
	}

	if film.Year != 0 {
		part = translator.Translate(session.Lang, "year", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>: %d\n", part, film.Year)
	}

	if film.Rating != 0 {
		part = translator.Translate(session.Lang, "rating", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>: %.2f★\n", part, film.Rating)
	}

	if film.IsViewed && film.UserRating != 0 {
		part = translator.Translate(session.Lang, "yourRating", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>: %.2f★\n", part, film.UserRating)
	}

	if film.Description != "" {
		part = translator.Translate(session.Lang, "description", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>:\n%s\n", part, film.Description)
	}

	if film.Comment != "" {
		part = translator.Translate(session.Lang, "comment", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>:\n%s\n", part, film.Comment)
	}

	if film.IsViewed && film.Review != "" {
		part = translator.Translate(session.Lang, "review", nil, nil)
		msg += fmt.Sprintf("<b>%s</b>:\n%s\n", part, film.Review)
	}

	return msg
}

func BuildFilmDetailWithNumberMessage(session *models.Session, itemID int, film *apiModels.Film) string {
	numberEmoji := utils.NumberToEmoji(itemID)

	msg := fmt.Sprintf("%s", numberEmoji)
	return msg + BuildFilmDetailMessage(session, film)
}

func BuildFilmGeneralMessage(session *models.Session, film *apiModels.Film) string {
	filmMsg := translator.Translate(session.Lang, "film", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>:  %s\n", filmMsg, film.Title)

	if film.Rating != 0 {
		ratingMsg := translator.Translate(session.Lang, "rating", nil, nil)
		msg += fmt.Sprintf("<b>%s:</b> %.2f★\n", ratingMsg, film.Rating)
	}

	if film.Genre != "" && film.Year != 0 {
		msg += "🎭 "

		if film.Genre != "" {
			msg += fmt.Sprintf("%d", film.Year)
		}

		if film.Genre != "" {
			if film.Year != 0 {
				msg += ", "
			}
			msg += fmt.Sprintf("%s", film.Genre)
		}

		msg += "\n"
	}

	if film.Description != "" {
		if len(film.Description) > 400 {
			film.Description, _ = utils.SplitTextByLength(film.Description, 300)
		}

		descriptionMsg := translator.Translate(session.Lang, "description", nil, nil)
		msg += fmt.Sprintf("<b>%s:</b> %s\n", descriptionMsg, film.Description)
	}

	msg += fmt.Sprintf("%s\n\n", boolToString(session, film.IsViewed))

	return msg
}

func boolToString(session *models.Session, viewed bool) string {
	if viewed {
		viewedMsg := translator.Translate(session.Lang, "viewed", nil, nil)
		msg := fmt.Sprintf("<b>%s</b>✔️", viewedMsg)
		return msg
	}

	notViewedMsg := translator.Translate(session.Lang, "notViewed", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>✖️", notViewedMsg)
	return msg
}
