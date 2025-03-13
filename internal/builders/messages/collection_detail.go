package messages

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

func CollectionDetail(collection *apiModels.Collection) string {
	return fmt.Sprintf("%s (%d)\n%s\n",
		toBold(collection.Name),
		collection.TotalFilms,
		formatOptionalString("", toItalic(collection.Description), "%s%s\n"))
}
