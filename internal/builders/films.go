package builders

//func BuildCollectionFilmsMessage(filmsResponse *models.FilmsResponse) string {
//	msg := "Вот ваши фильмы:\n"
//
//	for i, film := range filmsResponse.Films {
//		itemID := i + 1 + ((filmsResponse.Metadata.CurrentPage - 1) * filmsResponse.Metadata.PageSize)
//
//		msg += fmt.Sprintf("%d. ID: %d\nTitle: %s\n",
//			itemID, film.ID, film.Title)
//	}
//	msg += fmt.Sprintf("%d из %d страниц\n", filmsResponse.Metadata.CurrentPage, filmsResponse.Metadata.LastPage)
//
//	return msg
//}
//
//func BuildCollectionFilmsButtonCreate(filmResponse *models.FilmsResponse)
