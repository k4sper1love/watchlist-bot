package messages

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strings"
	"unicode/utf8"
)

func BuildFilmDetailMessage(session *models.Session, film *apiModels.Film) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("<b>%s</b>", film.Title))

	if film.Year != 0 {
		msg.WriteString(fmt.Sprintf(" <i>(%d)</i>", film.Year))
	}

	if film.IsFavorite {
		msg.WriteString(fmt.Sprintf(" ⭐"))
	}

	msg.WriteString("\n\n")

	details := make([]string, 0)
	details = append(details, utils.ViewedToEmojiColored(film.IsViewed))

	details = append(details, fmt.Sprintf("%d", film.ID))
	if film.Genre != "" {
		details = append(details, fmt.Sprintf("%s", film.Genre))
	}
	if film.Rating != 0 {
		details = append(details, fmt.Sprintf("★%.2f", film.Rating))
	}
	if film.IsViewed && film.UserRating != 0 {
		details = append(details, fmt.Sprintf("👤%.2f", film.UserRating))
	}
	if len(details) > 0 {
		msg.WriteString(strings.Join(details, " | ") + "\n\n")
	}

	if film.Description != "" {
		msg.WriteString(fmt.Sprintf("<b>Описание:</b>\n<i>%s</i>\n\n", film.Description))
	}

	if film.Comment != "" {
		commentMsg := translator.Translate(session.Lang, "comment", nil, nil)
		msg.WriteString(fmt.Sprintf("<b>%s:</b>\n<i>%s</i>\n\n", commentMsg, film.Comment))
	}

	if film.IsViewed && film.Review != "" {
		reviewMsg := translator.Translate(session.Lang, "review", nil, nil)
		msg.WriteString(fmt.Sprintf("<b>%s:</b>\n<i>%s</i>\n\n", reviewMsg, film.Review))
	}

	if session.Context == states.ContextCollection {
		msg.WriteString(fmt.Sprintf("📚 <i>%s</i>\n\n", session.CollectionDetailState.Collection.Name))
	}

	return msg.String()
}

func BuildFilmDetailWithNumberMessage(session *models.Session, itemID int, film *apiModels.Film) string {
	numberEmoji := utils.NumberToEmoji(itemID)

	msg := fmt.Sprintf("%s", numberEmoji)
	return msg + BuildFilmDetailMessage(session, film)
}

func BuildFilmGeneralMessage(session *models.Session, film *apiModels.Film, needViewed bool) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("<b>%s</b>", film.Title))

	if film.Year != 0 {
		msg.WriteString(fmt.Sprintf(" <i>(%d)</i>", film.Year))
	}

	msg.WriteString(fmt.Sprintf(" | %s\n", utils.ViewedToEmojiColored(film.IsViewed)))

	details := make([]string, 0)
	if film.Genre != "" {
		details = append(details, fmt.Sprintf("%s", film.Genre))
	}
	if film.Rating != 0 {
		details = append(details, fmt.Sprintf("★%.2f", film.Rating))
	}
	if needViewed && film.IsViewed && film.UserRating != 0 {
		details = append(details, fmt.Sprintf("👤%.2f", film.UserRating))
	}
	if len(details) > 0 {
		msg.WriteString(strings.Join(details, " , ") + "\n")
	}

	if film.Description != "" {
		if film.Genre == "YouTube Video" {
			parts := strings.Split(film.Description, "\n")
			parts = strings.Split(parts[0], ":")
			author := strings.TrimSpace(parts[1])
			film.Description = fmt.Sprintf("👨‍💼 <i>%s</i>", author)

			msg.WriteString(fmt.Sprintf("%s\n", film.Description))

		} else {
			if utf8.RuneCountInString(film.Description) > 230 {
				film.Description, _ = utils.SplitTextByLength(film.Description, 230)
			}

			descriptionMsg := translator.Translate(session.Lang, "description", nil, nil)
			msg.WriteString(fmt.Sprintf("<b>%s:</b>\n<i>%s</i>\n", descriptionMsg, film.Description))
		}
	}

	msg.WriteString("\n")

	return msg.String()
}

func viewedToString(session *models.Session, viewed bool) string {
	var messageCode string

	if viewed {
		messageCode = "viewed"
	} else {
		messageCode = "notViewed"
	}

	return translator.Translate(session.Lang, messageCode, nil, nil)
}
