package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func SaveFeedbackToDatabase(telegramID int, category, feedback string) error {
	feedbackEntry := models.Feedback{
		TelegramID: telegramID,
		Category:   category,
		Message:    feedback,
	}

	return GetDB().Create(&feedbackEntry).Error
}

func FetchAllFeedbacks() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := GetDB().Order("created_at DESC").Find(&feedbacks).Error
	return feedbacks, err
}

func DeleteFeedbackByID(id int) error {
	return GetDB().Delete(&models.Feedback{}, id).Error
}
