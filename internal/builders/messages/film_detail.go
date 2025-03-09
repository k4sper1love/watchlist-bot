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

func FilmDetail(session *models.Session) string {
	film := session.FilmDetailState.Film
	return fmt.Sprintf("%s%s%s\n\n%s%s%s%s%s",
		toBold(film.Title),
		toItalic(formatOptionalNumber("", film.Year, 0, "%s (%d)")),
		formatOptionalBool("⭐", film.IsFavorite, " %s"),
		formatFilmDetails(&film),
		formatOptionalString(toBold(translator.Translate(session.Lang, "description", nil, nil)),
			toItalic(film.Description), "%s:\n%s\n\n"),
		formatOptionalString(toBold(translator.Translate(session.Lang, "comment", nil, nil)),
			toItalic(film.Comment), "%s:\n%s\n\n"),
		formatOptionalString(toBold(translator.Translate(session.Lang, "review", nil, nil)),
			formatOptionalBool(toItalic(film.Review), film.IsViewed, "%s"), "%s:\n%s\n\n"),
		formatOptionalBool(toItalic(session.CollectionDetailState.Collection.Name), session.Context == states.ContextCollection, "📚 %s\n\n"))
}

func FilmGeneral(session *models.Session, film *apiModels.Film, needViewed bool) string {
	return fmt.Sprintf("%s%s | %s\n%s%s\n",
		toBold(film.Title),
		toItalic(fmt.Sprintf(" (%d)", film.Year)),
		utils.ViewedToEmojiColored(film.IsViewed),
		formatFilmGeneralDetails(film, needViewed),
		formatFilmGeneralDescription(session, film))
}

func formatFilmDetails(film *apiModels.Film) string {
	var details []string

	details = append(details, utils.ViewedToEmojiColored(film.IsViewed))
	details = append(details, fmt.Sprintf("%d", film.ID))

	if film.Genre != "" {
		details = append(details, film.Genre)
	}
	if film.Rating != 0 {
		details = append(details, fmt.Sprintf("★%.2f", film.Rating))
	}
	if film.IsViewed && film.UserRating != 0 {
		details = append(details, fmt.Sprintf("👤%.2f", film.UserRating))
	}

	if len(details) > 0 {
		return strings.Join(details, " | ") + "\n\n"
	}
	return ""
}

func formatFilmGeneralDetails(film *apiModels.Film, needViewed bool) string {
	var details []string

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
		return strings.Join(details, " , ") + "\n"
	}
	return ""
}

func formatFilmGeneralDescription(session *models.Session, film *apiModels.Film) string {
	if film.Description == "" {
		return ""
	}

	if film.Genre == "YouTube Video" {
		parts := strings.Split(strings.Split(film.Description, "\n")[0], ":")
		film.Description = fmt.Sprintf("👨‍💼 %s", toItalic(strings.TrimSpace(parts[1])))
		return fmt.Sprintf("%s\n\n", film.Description)
	}

	if utf8.RuneCountInString(film.Description) > 230 {
		film.Description, _ = utils.SplitTextByLength(film.Description, 230)
	}

	return fmt.Sprintf("%s:\n%s\n",
		toBold(translator.Translate(session.Lang, "description", nil, nil)),
		toItalic(film.Description),
	)
}
