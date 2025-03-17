package models

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

// BaseResponse represents a generic response structure containing metadata.
type BaseResponse struct {
	Metadata filters.Metadata `json:"metadata"` // Metadata includes pagination details.
}

// CollectionsResponse represents a response containing a list of collections.
type CollectionsResponse struct {
	BaseResponse                        // Embedded base response with metadata.
	Collections  []apiModels.Collection `json:"collections"` // List of collections returned by the API.
}

// FilmsResponse represents a response containing a list of films.
type FilmsResponse struct {
	BaseResponse                  // Embedded base response with metadata.
	Films        []apiModels.Film `json:"films"` // List of films returned by the API.
}

// CollectionFilmsResponse represents a response containing films associated with a specific collection.
type CollectionFilmsResponse struct {
	BaseResponse                              // Embedded base response with metadata.
	CollectionFilms apiModels.CollectionFilms `json:"collection_films"` // Films associated with a collection.
}

// CollectionFilmResponse represents a response containing a single collection and a single film.
type CollectionFilmResponse struct {
	Collection apiModels.Collection `json:"collection"` // The collection associated with the film.
	Film       apiModels.Film       `json:"film"`       // The film associated with the collection.
}
