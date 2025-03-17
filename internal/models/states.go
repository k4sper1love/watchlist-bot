package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"gorm.io/gorm"
)

// BaseState represents a base structure for all state types.
type BaseState struct {
	gorm.Model `json:"-"` // Embedded GORM model for database operations.
	SessionID  uint       `json:"-"`
}

// ProfileState represents the state for managing user profile details.
type ProfileState struct {
	BaseState
	Username string `json:"username,omitempty"` // User's username (optional).
	Email    string `json:"email,omitempty"`    // User's email (optional).
}

// FeedbackState represents the state for managing feedback submissions.
type FeedbackState struct {
	BaseState
	Category string // Category of the feedback (e.g., suggestions, bugs, issues).
	Message  string // Content of the feedback message.
}

// FilmsState represents the state for managing films and their filters/sorting.
type FilmsState struct {
	BaseState
	Films             []apiModels.Film `json:"films" gorm:"serializer:json"`                              // List of films in the current state.
	LastPage          int              `json:"-"`                                                         // Last page number for pagination.
	PageSize          int              `json:"-" gorm:"default:4"`                                        // Number of films per page.
	CurrentPage       int              `json:"-"`                                                         // Current page number.
	TotalRecords      int              `json:"-"`                                                         // Total number of films.
	Title             string           `json:"-"`                                                         // Search title for filtering films.
	FilmFilters       *FilmFilters     `gorm:"polymorphic:Filterable;polymorphicValue:FilmFilters"`       // Filters for films.
	CollectionFilters *FilmFilters     `gorm:"polymorphic:Filterable;polymorphicValue:CollectionFilters"` // Filters for collections.
	FilmSorting       *Sorting         `gorm:"polymorphic:Sortable;polymorphicValue:FilmSorting"`         // Sorting options for films.
	CollectionSorting *Sorting         `gorm:"polymorphic:Sortable;polymorphicValue:CollectionSorting"`   // Sorting options for collections.
}

// FilmDetailState represents the state for managing detailed information about a specific film.
type FilmDetailState struct {
	BaseState
	Index       int            `json:"-"`                           // Index of the film in the current list.
	Film        apiModels.Film `json:"film" gorm:"serializer:json"` // Detailed film data.
	IsFavorite  *bool          `json:"is_favorite,omitempty"`       // Indicates if the film is marked as a favorite.
	Title       string         `json:"title,omitempty"`             // Title of the film.
	Year        int            `json:"year,omitempty"`              // Release year of the film.
	Genre       string         `json:"genre,omitempty"`             // Genre of the film.
	Description string         `json:"description,omitempty"`       // Description of the film.
	Rating      float64        `json:"rating,omitempty"`            // Rating of the film.
	ImageURL    string         `json:"image_url,omitempty"`         // URL of the film's image.
	Comment     string         `json:"comment,omitempty"`           // User's comment on the film.
	IsViewed    *bool          `json:"is_viewed,omitempty"`         // Indicates if the film has been viewed.
	UserRating  float64        `json:"user_rating"`                 // User's rating for the film.
	Review      string         `json:"review"`                      // User's review for the film.
	URL         string         `json:"url,omitempty"`               // URL of the film.
}

// CollectionsState represents the state for managing collections and their sorting.
type CollectionsState struct {
	BaseState
	Collections []apiModels.Collection `json:"collections" gorm:"serializer:json"`            // List of collections in the current state.
	LastPage    int                    `json:"-"`                                             // Last page number for pagination.
	PageSize    int                    `json:"-" gorm:"default:4"`                            // Number of collections per page.
	CurrentPage int                    `json:"-"`                                             // Current page number.
	Name        string                 `json:"-"`                                             // Search name for filtering collections.
	Sorting     *Sorting               `gorm:"polymorphic:Sortable;polymorphicValue:Sorting"` // Sorting options for collections.
}

// CollectionDetailState represents the state for managing detailed information about a specific collection.
type CollectionDetailState struct {
	BaseState
	ObjectID    int                  `json:"-"`                                 // ID of the collection object.
	Collection  apiModels.Collection `json:"collection" gorm:"serializer:json"` // Detailed collection data.
	IsFavorite  *bool                `json:"is_favorite,omitempty"`             // Indicates if the collection is marked as a favorite.
	Name        string               `json:"name,omitempty"`                    // Name of the collection.
	Description string               `json:"description,omitempty"`             // Description of the collection.
}

// CollectionFilmsState represents the state for managing films within a collection.
type CollectionFilmsState struct {
	BaseState
	LastPage     int `json:"-"`                  // Last page number for pagination.
	PageSize     int `json:"-" gorm:"default:4"` // Number of films per page.
	CurrentPage  int `json:"-"`                  // Current page number.
	TotalRecords int `json:"-"`                  // Total number of films in the collection.
}

// AdminState represents the state for managing admin-related tasks.
type AdminState struct {
	BaseState
	IsAdmin      bool       `json:"-"`                  // Indicates if the user is an admin.
	UserID       int        `json:"-"`                  // ID of the user being managed.
	UserLang     string     `json:"-"`                  // Language of the user being managed.
	UserRole     roles.Role `json:"-"`                  // Role of the user being managed.
	FeedbackID   int        `json:"-"`                  // ID of the feedback being managed.
	LastPage     int        `json:"-"`                  // Last page number for pagination.
	PageSize     int        `json:"-" gorm:"default:4"` // Number of items per page.
	CurrentPage  int        `json:"-"`                  // Current page number.
	TotalRecords int        `json:"-"`                  // Total number of items.
	Message      string     `json:"-"`                  // Message content for admin tasks.
	ImageURL     string     `json:"-"`                  // URL of the image for admin tasks.
	NeedPin      bool       `json:"-"`                  // Indicates if the message needs to be pinned.
}

// Clear resets the state of films, including the title and sorting options.
func (s *FilmsState) Clear() {
	s.Title = ""
	s.FilmSorting.Clear()
	s.CollectionSorting.Clear()
}

// Clear resets the state of collections, including the name and sorting options.
func (s *CollectionsState) Clear() {
	s.Name = ""
	s.Sorting.Clear()
}

// Clear resets the state of admin tasks, including the message and image URL.
func (s *AdminState) Clear() {
	s.Message, s.ImageURL = "", ""
	s.NeedPin = false
}

// ResetAdmin resets the admin-specific fields (e.g., IsAdmin).
func (s *AdminState) ResetAdmin() {
	s.IsAdmin = false
}

// Clear resets the state of profile details, including the username and email.
func (s *ProfileState) Clear() {
	s.Username, s.Email = "", ""
}

// Clear resets the state of feedback submissions, including the category and message.
func (s *FeedbackState) Clear() {
	s.Category, s.Message = "", ""
}

// Clear resets the state of film details, including all fields.
func (s *FilmDetailState) Clear() {
	s.Title = ""
	s.Year = 0
	s.Genre = ""
	s.Description = ""
	s.Rating = 0
	s.ImageURL = ""
	s.Comment = ""
	s.IsViewed = nil
	s.IsFavorite = nil
	s.UserRating = 0
	s.Review = ""
	s.URL = ""
}

// UpdateFilm updates the film state with new film data and clears the index.
func (s *FilmDetailState) UpdateFilm(film apiModels.Film) {
	s.Film = film
	s.ClearIndex()
}

// SetFavorite sets the favorite status of the film.
func (s *FilmDetailState) SetFavorite(value bool) {
	s.IsFavorite = &value
}

// SetViewed sets the viewed status of the film.
func (s *FilmDetailState) SetViewed(value bool) {
	s.IsViewed = &value
}

// IsViewedEdit checks if the viewed status has been edited.
func (s *FilmDetailState) IsViewedEdit() bool {
	return s.IsViewed != nil
}

// SyncValues synchronizes film values with default values if not explicitly set.
func (s *FilmDetailState) SyncValues() {
	if s.IsViewedEdit() {
		return
	}

	if s.UserRating == 0 {
		s.UserRating = s.Film.UserRating
	}
	if s.Review == "" {
		s.Review = s.Film.Review
	}
}

// HasIndex checks if the film index is set.
func (s *FilmDetailState) HasIndex() bool {
	return s.Index != -1
}

// ClearIndex resets the film index.
func (s *FilmDetailState) ClearIndex() {
	s.Index = -1
}

// Clear resets the state of collection details, including the favorite status, name, and description.
func (s *CollectionDetailState) Clear() {
	s.IsFavorite = nil
	s.Name, s.Description = "", ""
}

// SetImageURL sets the image URL for the film.
func (s *FilmDetailState) SetImageURL(url string) {
	s.ImageURL = url
}

// SetFromFilm populates the film detail state from a film object.
func (s *FilmDetailState) SetFromFilm(film *apiModels.Film) {
	s.Title = film.Title
	s.Description = film.Description
	s.Genre = film.Genre
	s.Year = film.Year
	s.Rating = film.Rating
	s.URL = film.URL
	s.ImageURL = film.ImageURL
}

// SetFavorite sets the favorite status of the collection.
func (s *CollectionDetailState) SetFavorite(value bool) {
	s.IsFavorite = &value
}
