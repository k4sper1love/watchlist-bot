package messages

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildFilmDetailMessage(film *apiModels.Film) string {
	msg := ""

	if film.Title != "" {
		msg += fmt.Sprintf("<b>Название</b>: %s\n", film.Title)
	}

	msg += fmt.Sprintf("<b>Просмотрено</b>: %s\n", boolToEmoji(film.IsViewed))

	if film.Genre != "" {
		msg += fmt.Sprintf("<b>Жанр</b>: %s\n", film.Genre)
	}

	if film.Year != 0 {
		msg += fmt.Sprintf("<b>Год выпуска</b>: %d\n", film.Year)
	}

	if film.Rating != 0 {
		msg += fmt.Sprintf("<b>Рейтинг</b>: %.2f★\n", film.Rating)
	}

	if film.IsViewed && film.UserRating != 0 {
		msg += fmt.Sprintf("<b>Ваша оценка</b>: %.2f★\n", film.UserRating)
	}

	if film.Description != "" {
		msg += fmt.Sprintf("<b>Описание</b>:\n%s\n", film.Description)
	}

	if film.Comment != "" {
		msg += fmt.Sprintf("<b>Комментарий</b>:\n%s\n", film.Comment)
	}

	if film.IsViewed && film.Review != "" {
		msg += fmt.Sprintf("<b>Рецензия</b>:\n%s\n", film.Review)
	}

	return msg
}

func BuildFilmDetailWithNumberMessage(itemID int, film *apiModels.Film) string {
	numberEmoji := numberToEmoji(itemID)

	msg := fmt.Sprintf("%s\n", numberEmoji)
	return msg + BuildFilmDetailMessage(film)
}

func BuildFilmGeneralMessage(film *apiModels.Film) string {
	msg := fmt.Sprintf("<b>Фильм</b>:  %s\n", film.Title)

	if film.Rating != 0 {
		msg += fmt.Sprintf("<b>Рейтинг:</b> %.2f★\n", film.Rating)
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
			fmt.Println(film.Description)
		}

		msg += fmt.Sprintf("<b>Описание:</b> %s\n", film.Description)
	}

	msg += fmt.Sprintf("%s\n\n", boolToString(film.IsViewed))

	return msg
}

func boolToEmoji(viewed bool) string {
	if viewed {
		return "✔️"
	}
	return "✖️"
}

func boolToString(viewed bool) string {
	if viewed {
		return "<b>Просмотрено </b>✔️"
	}
	return "<b>Не просмотрено </b>✖️"
}
