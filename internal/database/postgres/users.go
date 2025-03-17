package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

const (
	TelegramIDField       = "telegram_id"       // Field name for Telegram ID in the database.
	TelegramUsernameField = "telegram_username" // Field name for Telegram username in the database.
)

// GetUserCount returns the total number of users, optionally filtered by admin status.
func GetUserCount(isAdmin bool) (int64, error) {
	var count int64

	query := GetDatabase().Model(&models.Session{})
	if isAdmin {
		query = query.Where("role > 0")
	}

	err := query.Count(&count).Error
	return count, err
}

// GetTelegramIDs retrieves all Telegram IDs from the sessions table.
func GetTelegramIDs() ([]int, error) {
	var telegramIDs []int
	err := GetDatabase().Model(&models.Session{}).Pluck("telegram_id", &telegramIDs).Error
	return telegramIDs, err
}

// GetUsers retrieves all users, optionally filtered by admin status, ordered by creation time.
func GetUsers(isAdmin bool) ([]models.Session, error) {
	var sessions []models.Session

	query := GetDatabase().Order("created_at DESC")
	if isAdmin {
		query = query.Where("role > 0")
	}

	err := query.Find(&sessions).Error
	return sessions, err
}

// GetUsersWithPagination retrieves users with pagination support, optionally filtered by admin status.
func GetUsersWithPagination(page, pageSize int, isAdmin bool) ([]models.Session, error) {
	var sessions []models.Session
	offset := utils.CalculateOffset(page, pageSize)

	query := GetDatabase().Order("created_at DESC").Limit(pageSize).Offset(offset)
	if isAdmin {
		query = query.Where("role > 0")
	}

	err := query.Find(&sessions).Error
	return sessions, err
}

// SetUserBanStatus updates the ban status of a user by Telegram ID.
func SetUserBanStatus(telegramID int, isBanned bool) error {
	return GetDatabase().Model(&models.Session{}).Where("telegram_id = ?", telegramID).Update("is_banned", isBanned).Error
}

// IsUserBanned checks if a user is banned by Telegram ID.
func IsUserBanned(telegramID int) bool {
	var session models.Session
	if err := GetDatabase().Model(&models.Session{}).Where("telegram_id = ?", telegramID).First(&session).Error; err != nil {
		return false
	}
	return session.IsBanned
}

// SetUserRole updates the role of a user by Telegram ID.
func SetUserRole(telegramID int, role roles.Role) error {
	return GetDatabase().Model(&models.Session{}).Where("telegram_id = ?", telegramID).Update("role", role).Error
}

// GetUserByField retrieves a user by a specific field value, optionally filtered by admin status.
func GetUserByField(field string, value any, isAdmin bool) (*models.Session, error) {
	var session models.Session

	query := GetDatabase().Model(&models.Session{}).Where(field+" = ?", value)
	if isAdmin {
		query = query.Where("role > 0")
	}

	err := query.First(&session).Error
	return &session, err
}

// GetUserByAPIUserID retrieves a user session by API user ID, optionally filtered by admin status.
func GetUserByAPIUserID(id int, isAdmin bool) (*models.Session, error) {
	sessions, err := GetUsers(isAdmin)
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if session.User.ID == id {
			return &session, nil
		}
	}
	return nil, errors.New("session not found")
}
