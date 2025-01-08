package models

import (
	"gorm.io/gorm"
	"strings"
)

type Feedback struct {
	gorm.Model
	TelegramID int    `gorm:"not null"`
	Category   string `gorm:"not null"`
	Message    string `gorm:"not null"`
}

type FilmsFilters struct {
	gorm.Model
	FilterableID   uint   `json:"-"` // ID связанной сущности
	FilterableType string `json:"-"` // Тип связанной сущности (FilmFilters или CollectionFilters)
	MinRating      float64
	MaxRating      float64
}

type FilmsSorting struct {
	gorm.Model
	SortableID   uint   `json:"-"` // ID связанной сущности
	SortableType string `json:"-"` // Тип связанной сущности (FilmSorting или CollectionSorting)
	Field        string `json:"-"` // Поле для сортировки (например, "rating")
	Direction    string `json:"-"` // Направление ("asc" или "desc")
	Sort         string `json:"-"` // Любая дополнительная информация о сортировке
}

func (f *FilmsSorting) Clear() {
	f.Field = ""
	f.Direction = ""
}

func (f *FilmsFilters) IsFiltersEnabled() bool {
	if f.MaxRating > 0 || f.MinRating > 0 {
		return true
	}

	return false
}

func (f *FilmsFilters) ResetFilters() {
	f.MinRating = 0
	f.MaxRating = 0
}

func (f *FilmsFilters) IsFilterEnabled(filterType string) bool {
	switch filterType {
	case "minRating":
		return f.MinRating != 0
	case "maxRating":
		return f.MaxRating != 0
	default:
		return false
	}
}

func (f *FilmsSorting) IsSortingEnabled() bool {
	return f.Sort != ""
}

func (f *FilmsSorting) IsSortingFieldEnabled(field string) bool {
	return field == strings.TrimPrefix(f.Sort, "-")
}

func (f *FilmsSorting) ResetSorting() {
	f.Sort = ""
}
