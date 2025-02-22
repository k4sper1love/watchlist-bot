package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func SaveFeedback(telegramID int, telegramUsername, category, message string) error {
	entry := models.Feedback{
		TelegramID:       telegramID,
		TelegramUsername: telegramUsername,
		Category:         category,
		Message:          message,
	}

	return GetDatabase().Create(&entry).Error
}

func GetFeedbackCount() (int64, error) {
	var count int64
	err := GetDatabase().Model(&models.Feedback{}).Count(&count).Error
	return count, err
}

func GetFeedbacks() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := GetDatabase().Order("created_at DESC").Find(&feedbacks).Error
	return feedbacks, err
}

func GetFeedbacksWithPagination(page, pageSize int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	offset := utils.CalculateOffset(page, pageSize)
	err := GetDatabase().Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&feedbacks).Error
	return feedbacks, err
}

func GetFeedbackByID(id int) (*models.Feedback, error) {
	var feedback models.Feedback
	err := GetDatabase().Model(&models.Feedback{}).Where("id = ?", id).First(&feedback).Error
	return &feedback, err
}

func DeleteFeedbackByID(id int) error {
	return GetDatabase().Delete(&models.Feedback{}, id).Error
}
