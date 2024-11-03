package models

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

type CollectionsResponse struct {
	Collections []apiModels.Collection `json:"collections"`
	Metadata    filters.Metadata       `json:"metadata"`
}

type FilmsResponse struct {
	Films    []apiModels.Film `json:"films"`
	Metadata filters.Metadata `json:"metadata"`
}

type CollectionFilmsResponse struct {
	CollectionFilms apiModels.CollectionFilms `json:"collection_films"`
	Metadata        filters.Metadata          `json:"metadata"`
}

type CollectionFilmResponse struct {
	Collection apiModels.Collection `json:"collection"`
	Film       apiModels.Film       `json:"film"`
}
