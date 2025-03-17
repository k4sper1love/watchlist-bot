package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// SaveFeedback saves a new feedback entry to the database.
func SaveFeedback(telegramID int, telegramUsername, category, message string) error {
	entry := models.Feedback{
		TelegramID:       telegramID,
		TelegramUsername: telegramUsername,
		Category:         category,
		Message:          message,
	}

	return GetDatabase().Create(&entry).Error
}

// GetFeedbackCount returns the total number of feedback entries in the database.
func GetFeedbackCount() (int64, error) {
	var count int64
	err := GetDatabase().Model(&models.Feedback{}).Count(&count).Error
	return count, err
}

// GetFeedbacks retrieves all feedback entries ordered by creation time (newest first).
func GetFeedbacks() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := GetDatabase().Order("created_at DESC").Find(&feedbacks).Error
	return feedbacks, err
}

// GetFeedbacksWithPagination retrieves feedback entries with pagination support.
func GetFeedbacksWithPagination(page, pageSize int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	offset := utils.CalculateOffset(page, pageSize)
	err := GetDatabase().Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&feedbacks).Error
	return feedbacks, err
}

// GetFeedbackByID retrieves a feedback entry by its unique ID.
func GetFeedbackByID(id int) (*models.Feedback, error) {
	var feedback models.Feedback
	err := GetDatabase().Model(&models.Feedback{}).Where("id = ?", id).First(&feedback).Error
	return &feedback, err
}

// DeleteFeedbackByID deletes a feedback entry by its unique ID.
func DeleteFeedbackByID(id int) error {
	return GetDatabase().Delete(&models.Feedback{}, id).Error
}
