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

type FiltersFilm struct {
	gorm.Model
	FilterableID   uint   `json:"-"`
	FilterableType string `json:"-"`
	MinRating      float64
	MaxRating      float64
}

type Sorting struct {
	gorm.Model
	SortableID   uint   `json:"-"`
	SortableType string `json:"-"`
	Field        string `json:"-"`
	Direction    string `json:"-"`
	Sort         string `json:"-"`
}

func (f *Sorting) Clear() {
	f.Field = ""
	f.Direction = ""
}

func (f *FiltersFilm) IsFiltersEnabled() bool {
	if f.MaxRating > 0 || f.MinRating > 0 {
		return true
	}

	return false
}

func (f *FiltersFilm) ResetFilters() {
	f.MinRating = 0
	f.MaxRating = 0
}

func (f *FiltersFilm) IsFilterEnabled(filterType string) bool {
	switch filterType {
	case "minRating":
		return f.MinRating != 0
	case "maxRating":
		return f.MaxRating != 0
	default:
		return false
	}
}

func (f *Sorting) IsSortingEnabled() bool {
	return f.Sort != ""
}

func (f *Sorting) IsSortingFieldEnabled(field string) bool {
	return field == strings.TrimPrefix(f.Sort, "-")
}

func (f *Sorting) ResetSorting() {
	f.Sort = ""
}
