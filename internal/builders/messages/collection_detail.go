package messages

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildCollectionDetailMessage(session *models.Session, collection *apiModels.Collection) string {
	collectionMsg := translator.Translate(session.Lang, "collection", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b> %s\n", collectionMsg, collection.Name)

	if collection.Description != "" {
		descriptionMsg := translator.Translate(session.Lang, "description", nil, nil)
		msg += fmt.Sprintf("<b>%s:</b> %s\n", descriptionMsg, collection.Description)
	}

	totalFilmsMsg := translator.Translate(session.Lang, "totalFilms", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %d\n\n", totalFilmsMsg, collection.TotalFilms)
	return msg
}
