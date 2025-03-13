package models

import (
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Feedback struct {
	gorm.Model
	TelegramID       int `gorm:"not null"`
	TelegramUsername string
	Category         string `gorm:"not null"`
	Message          string `gorm:"not null"`
}

type FilmFilters struct {
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

func (f *FilmFilters) ResetAll() {
	f.Rating = ""
	f.UserRating = ""
	f.Year = ""
	f.IsViewed = nil
	f.IsFavorite = nil
	f.HasURL = nil
}

func (f *FilmFilters) Reset(filterType string) {
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

func (f *FilmFilters) IsEnabled() bool {
	return f.Rating != "" || f.UserRating != "" || f.Year != "" || f.IsViewed != nil || f.IsFavorite != nil || f.HasURL != nil
}

func (f *FilmFilters) IsFieldEnabled(filterType string) bool {
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

func (f *FilmFilters) ApplyRange(filterType, value string) {
	switch filterType {
	case "rating":
		f.Rating = value
	case "userRating":
		f.UserRating = value
	case "year":
		f.Year = value
	}
}

func (f *FilmFilters) ApplySwitch(filterType string, value bool) {
	switch filterType {
	case "isViewed":
		f.IsViewed = &value
	case "isFavorite":
		f.IsFavorite = &value
	case "hasURL":
		f.HasURL = &value
	}
}

func (f *FilmFilters) ToString(filterType string) string {
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

func (f *Sorting) Reset() {
	f.Sort = ""
}

func (f *Sorting) IsEnabled() bool {
	return f.Sort != ""
}

func (f *Sorting) IsFieldEnabled(field string) bool {
	return field == strings.TrimPrefix(f.Sort, "-")
}

func (f *Sorting) SetSort() {
	f.Sort = f.Direction + f.Sort
}
