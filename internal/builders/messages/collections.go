package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildCollectionsMessage(session *models.Session, metadata *filters.Metadata) string {
	collections := session.CollectionsState.Collections

	msg := ""

	if metadata.TotalRecords == 0 {
		msg += "Не найдено коллекций."
		return msg
	}

	msg += fmt.Sprintf("📚 <b>Всего коллекций:</b> %d\n\n", metadata.TotalRecords)

	for i, collection := range collections {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)

		numberEmoji := numberToEmoji(itemID)

		msg += fmt.Sprintf("%s\n", numberEmoji)
		msg += BuildCollectionDetailMessage(&collection)
	}

	msg += fmt.Sprintf("<b>📄 Страница %d из %d</b>\n\n", metadata.CurrentPage, metadata.LastPage)
	msg += "Выберите коллекцию из списка, чтобы узнать больше."

	return msg
}

func numberToEmoji(number int) string {
	emojis := []string{"0️⃣", "1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣"}
	if number < 10 {
		return emojis[number]
	}

	result := ""
	for number > 0 {
		digit := number % 10
		result = emojis[digit] + result
		number /= 10
	}
	return result
}
