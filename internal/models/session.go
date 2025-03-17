package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"gorm.io/gorm"
)

// Session represents the session state of a user in the Watchlist bot.
type Session struct {
	gorm.Model                                   // Embedded GORM model for database operations.
	TelegramID            int                    `gorm:"unique"` // Unique Telegram user ID.
	TelegramUsername      string                 // Optional Telegram username.
	Role                  roles.Role             `json:"role" gorm:"serializer:json"` // User's role (e.g., user, admin, super admin).
	Lang                  string                 // User's preferred language.
	IsBanned              bool                   `gorm:"default:false"`               // Indicates if the user is banned.
	User                  apiModels.User         `json:"user" gorm:"serializer:json"` // Associated API user data.
	AccessToken           string                 `json:"access_token"`                // Access token for API authentication.
	RefreshToken          string                 `json:"refresh_token"`               // Refresh token for API authentication.
	KinopoiskAPIToken     string                 `json:"kinopoisk_api_token"`         // Kinopoisk API token for external API requests.
	State                 string                 // Current session state (e.g., awaiting input).
	Context               string                 // Current session context (e.g., film, collection).
	AdminState            *AdminState            `gorm:"foreignKey:SessionID"` // Admin-specific session state.
	ProfileState          *ProfileState          `gorm:"foreignKey:SessionID"` // Profile-specific session state.
	FeedbackState         *FeedbackState         `gorm:"foreignKey:SessionID"` // Feedback-specific session state.
	CollectionsState      *CollectionsState      `gorm:"foreignKey:SessionID"` // Collections-specific session state.
	CollectionDetailState *CollectionDetailState `gorm:"foreignKey:SessionID"` // Collection detail-specific session state.
	FilmsState            *FilmsState            `gorm:"foreignKey:SessionID"` // Films-specific session state.
	FilmDetailState       *FilmDetailState       `gorm:"foreignKey:SessionID"` // Film detail-specific session state.
	CollectionFilmsState  *CollectionFilmsState  `gorm:"foreignKey:SessionID"` // Collection films-specific session state.
}

// SetContext sets the current session context (e.g., film or collection).
func (s *Session) SetContext(context string) {
	s.Context = context
}

// SetState sets the current session state (e.g., awaiting input).
func (s *Session) SetState(state string) {
	s.State = state
}

// ClearState clears the current session state.
func (s *Session) ClearState() {
	s.State = ""
}

// ClearContext clears the current session context.
func (s *Session) ClearContext() {
	s.Context = ""
}

// ClearUser clears the associated API user data.
func (s *Session) ClearUser() {
	s.User = apiModels.User{}
}

// ClearAllStates resets all session states to their default values.
func (s *Session) ClearAllStates() {
	s.ClearState()
	s.ProfileState.Clear()
	s.FeedbackState.Clear()
	s.FilmsState.Clear()
	s.FilmDetailState.Clear()
	s.CollectionsState.Clear()
	s.CollectionDetailState.Clear()
	s.AdminState.Clear()
}

// Logout logs the user out by clearing tokens, user data, context, and all states.
func (s *Session) Logout() {
	s.AccessToken, s.RefreshToken, s.KinopoiskAPIToken = "", "", ""
	s.ClearUser()
	s.ClearContext()
	s.ClearAllStates()
}

// GetFilmFiltersByCtx retrieves the film filters based on the current session context.
func (s *Session) GetFilmFiltersByCtx() *FilmFilters {
	switch s.Context {
	case states.CtxFilm:
		return s.FilmsState.FilmFilters
	case states.CtxCollection:
		return s.FilmsState.CollectionFilters
	default:
		return &FilmFilters{}
	}
}

// GetFilmSortingByCtx retrieves the film sorting options based on the current session context.
func (s *Session) GetFilmSortingByCtx() *Sorting {
	switch s.Context {
	case states.CtxFilm:
		return s.FilmsState.FilmSorting
	case states.CtxCollection:
		return s.FilmsState.CollectionSorting
	default:
		return &Sorting{}
	}
}
