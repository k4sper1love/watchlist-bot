package postgres

import "github.com/k4sper1love/watchlist-bot/internal/models"

func GetSessionByTelegramID(telegramID int) (*models.Session, error) {
	var session models.Session

	if err := GetDB().FirstOrInit(&session, models.Session{TelegramID: telegramID}).Error; err != nil {
		return nil, err
	}

	return &session, nil
}
