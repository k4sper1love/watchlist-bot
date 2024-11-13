package builders

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

func BuildFilmDetailMessage(film *apiModels.Film) string {
	msg := fmt.Sprintf("🆔 <b>ID</b>: %d\n", film.ID)

	if film.Title != "" {
		msg += fmt.Sprintf("📛 <b>Название</b>: %s\n", film.Title)
	}
	if film.Description != "" {
		msg += fmt.Sprintf("📝 <b>Описание</b>: %s\n", film.Description)
	}
	if film.Genre != "" {
		msg += fmt.Sprintf("🎭 <b>Жанр</b>: %s\n", film.Genre)
	}
	if film.Rating != 0 {
		msg += fmt.Sprintf("⭐️ <b>Рейтинг</b>: %.2f\n", film.Rating)
	}

	if film.Year != 0 {
		msg += fmt.Sprintf("📅 <b>Год выпуска</b>: %d\n", film.Year)
	}

	if film.Comment != "" {
		msg += fmt.Sprintf("💬 <b>Комментарий</b>: %s\n", film.Comment)
	}

	msg += fmt.Sprintf("📈 <b>Просмотрено</b>: %s\n", boolToString(film.IsViewed))

	if !film.IsViewed {
		return msg
	}

	if film.UserRating != 0 {
		msg += fmt.Sprintf("👤 <b>Оценка пользователя</b>: %.2f\n", film.UserRating)
	}

	if film.Review != "" {
		msg += fmt.Sprintf("🖋️ <b>Рецензия</b>: %s\n", film.Review)
	}

	return msg
}

func BuildFilmGeneralMessage(itemID int, film *apiModels.Film) string {
	msg := fmt.Sprintf("<b>🎬 Фильм №%d</b>. ", itemID)
	msg += fmt.Sprintf("<b>ID:</b> %d\n\n", film.ID)

	msg += fmt.Sprintf("🎥 <b>Название:</b> %s\n", film.Title)

	msg += fmt.Sprintf("📈 <b>Просмотрено</b>: %s\n", boolToString(film.IsViewed))

	if film.Genre != "" || film.Year != 0 || film.Rating != 0 {
		msg += "🎭 <b>Информация:</b> "

		if film.Genre != "" {
			msg += fmt.Sprintf("%s", film.Genre)
		}
		if film.Year != 0 {
			if film.Genre != "" {
				msg += ", "
			}
			msg += fmt.Sprintf("%d г.", film.Year)
		}
		if film.Rating != 0 {
			if film.Genre != "" || film.Year != 0 {
				msg += ", "
			}
			msg += fmt.Sprintf("%.2f ⭐️", film.Rating)
		}
		msg += "\n"
	}

	if film.Description != "" {
		msg += fmt.Sprintf("📝 <b>Описание:</b> %s\n", film.Description)
	}

	msg += fmt.Sprintf("\n🕒 <b>Последнее обновление:</b> %s\n", film.UpdatedAt.Format("02.01.06 15:04"))
	msg += fmt.Sprintf("📅 <b>Создан:</b> %s\n\n", film.CreatedAt.Format("02.01.06 15:04"))

	return msg
}

func boolToString(viewed bool) string {
	if viewed {
		return "✅"
	}
	return "❌"
}
