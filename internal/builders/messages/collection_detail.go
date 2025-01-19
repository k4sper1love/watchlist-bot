package messages

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func BuildCollectionDetailMessage(session *models.Session, collection *apiModels.Collection) string {
	msg := fmt.Sprintf(fmt.Sprintf("<b>%s</b>", collection.Name))

	msg += fmt.Sprintf(" (%d)\n", collection.TotalFilms)

	if collection.Description != "" {
		msg += fmt.Sprintf("<i>%s</i>\n", collection.Description)
	}

	msg += "\n"

	return msg
}
