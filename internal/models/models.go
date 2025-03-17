package models

import (
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// Feedback represents a feedback entry submitted by a user.
type Feedback struct {
	gorm.Model              // Embedded GORM model for database operations.
	TelegramID       int    `gorm:"not null"` // Telegram user ID of the feedback submitter.
	TelegramUsername string // Optional Telegram username of the feedback submitter.
	Category         string `gorm:"not null"` // Category of the feedback (e.g., suggestions, bugs, issues).
	Message          string `gorm:"not null"` // The content of the feedback message.
}

// FilmFilters represents filters applied to films for searching or sorting.
type FilmFilters struct {
	gorm.Model            // Embedded GORM model for database operations.
	FilterableID   uint   `json:"-"` // ID of the entity being filtered (e.g., film or collection).
	FilterableType string `json:"-"` // Type of the entity being filtered (e.g., "film", "collection").
	Rating         string `json:"-"` // Filter for film rating range (e.g., "7-9").
	UserRating     string `json:"-"` // Filter for user rating range (e.g., "3-5").
	Year           string `json:"-"` // Filter for film release year range (e.g., "2000-2010").
	IsViewed       *bool  `json:"-"` // Filter for whether the film has been viewed.
	IsFavorite     *bool  `json:"-"` // Filter for whether the film is marked as a favorite.
	HasURL         *bool  `json:"-"` // Filter for whether the film has an associated URL.
}

// Sorting represents sorting options applied to entities like films or collections.
type Sorting struct {
	gorm.Model          // Embedded GORM model for database operations.
	SortableID   uint   `json:"-"` // ID of the entity being sorted (e.g., film or collection).
	SortableType string `json:"-"` // Type of the entity being sorted (e.g., "film", "collection").
	Field        string `json:"-"` // Field to sort by (e.g., "title", "rating").
	Direction    string `json:"-"` // Sorting direction (empty for asc or "-" for desc).
	Sort         string `json:"-"` // Combined sort string (e.g., "-title").
}

// ResetAll resets all filters to their default values.
func (f *FilmFilters) ResetAll() {
	f.Rating = ""
	f.UserRating = ""
	f.Year = ""
	f.IsViewed = nil
	f.IsFavorite = nil
	f.HasURL = nil
}

// Reset resets a specific filter based on its type.
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

// IsEnabled checks if any filter is currently active.
func (f *FilmFilters) IsEnabled() bool {
	return f.Rating != "" || f.UserRating != "" || f.Year != "" || f.IsViewed != nil || f.IsFavorite != nil || f.HasURL != nil
}

// IsFieldEnabled checks if a specific filter field is active.
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

// ApplyRange applies a range-based filter (e.g., rating, user rating, year).
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

// ApplySwitch applies a switch-based filter (e.g., isViewed, isFavorite, hasURL).
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

// String converts a specific filter field to a string representation.
func (f *FilmFilters) String(filterType string) string {
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

// Clear resets all sorting fields to their default values.
func (f *Sorting) Clear() {
	f.Field = ""
	f.Direction = ""
}

// Reset resets the sorting field to its default value.
func (f *Sorting) Reset() {
	f.Sort = ""
}

// IsEnabled checks if sorting is currently active.
func (f *Sorting) IsEnabled() bool {
	return f.Sort != ""
}

// IsFieldEnabled checks if a specific field is being used for sorting.
func (f *Sorting) IsFieldEnabled(field string) bool {
	return field == strings.TrimPrefix(f.Sort, "-")
}

// SetSort constructs the combined sort string based on the field and direction.
func (f *Sorting) SetSort() {
	f.Sort = f.Direction + f.Field
}
