package models

import (
	"gorm.io/gorm"
	"strconv"
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
	Rating         string `json:"-"`
	UserRating     string `json:"-"`
	Year           string `json:"-"`
	IsViewed       *bool  `json:"-"`
	IsFavorite     *bool  `json:"-"`
	HasURL         *bool  `json:"-"`
}

type FilterRangeConfig struct {
	MinValue float64
	MaxValue float64
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
	if f.Rating != "" {
		return true
	} else if f.UserRating != "" {
		return true
	} else if f.Year != "" {
		return true
	} else if f.IsViewed != nil {
		return true
	} else if f.IsFavorite != nil {
		return true
	} else if f.HasURL != nil {
		return true
	}

	return false
}

func (f *FiltersFilm) ResetFilters() {
	f.Rating = ""
	f.UserRating = ""
	f.Year = ""
	f.IsViewed = nil
	f.IsFavorite = nil
	f.HasURL = nil
}

func (f *FiltersFilm) ResetFilter(filterType string) {
	switch filterType {
	case "rating":
		f.Rating = ""
	case "userRating":
		f.UserRating = ""
	case "year":
		f.Year = ""
	case "isViewed":
		f.IsViewed = nil
	case "isFavorite":
		f.IsFavorite = nil
	case "hasURL":
		f.HasURL = nil
	}
}

func (f *FiltersFilm) IsFilterEnabled(filterType string) bool {
	switch filterType {
	case "rating":
		return f.Rating != ""
	case "userRating":
		return f.UserRating != ""
	case "year":
		return f.Year != ""
	case "isViewed":
		return f.IsViewed != nil
	case "isFavorite":
		return f.IsFavorite != nil
	case "hasURL":
		return f.HasURL != nil
	default:
		return false
	}
}

func (f *FiltersFilm) ApplyRangeValue(filterType, value string) {
	switch filterType {
	case "rating":
		f.Rating = value
	case "userRating":
		f.UserRating = value
	case "year":
		f.Year = value
	}
}

func (f *FiltersFilm) ApplySwitchValue(filterType string, value bool) {
	switch filterType {
	case "isViewed":
		f.IsViewed = &value
	case "isFavorite":
		f.IsFavorite = &value
	case "hasURL":
		f.HasURL = &value
	}
}

func (f *FiltersFilm) ValueToString(filterType string) string {
	switch filterType {
	case "rating":
		return f.Rating
	case "userRating":
		return f.UserRating
	case "year":
		return f.Year
	case "isViewed":
		return strconv.FormatBool(*f.IsViewed)
	case "isFavorite":
		return strconv.FormatBool(*f.IsFavorite)
	case "hasURL":
		return strconv.FormatBool(*f.HasURL)
	default:
		return ""
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
