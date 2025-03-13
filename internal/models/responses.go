package models

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

type BaseResponse struct {
	Metadata filters.Metadata `json:"metadata"`
}

type CollectionsResponse struct {
	BaseResponse
	Collections []apiModels.Collection `json:"collections"`
}

type FilmsResponse struct {
	BaseResponse
	Films []apiModels.Film `json:"films"`
}

type CollectionFilmsResponse struct {
	BaseResponse
	CollectionFilms apiModels.CollectionFilms `json:"collection_films"`
}

type CollectionFilmResponse struct {
	Collection apiModels.Collection `json:"collection"`
	Film       apiModels.Film       `json:"film"`
}
