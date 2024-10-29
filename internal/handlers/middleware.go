package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
)

func requireAuth(app config.App, session *models.Session, next func(config.App, *models.Session)) {
	if session.UserID == -1 {
		sendMessage(app, "Для выполнения этой команды нужно быть авторизованным. Используйте /start")
		resetState(session)
		return
	}

	if !watchlist.IsTokenValid(app, session.AccessToken) {
		if err := watchlist.RefreshAccessToken(app, session); err != nil {
			sendMessage(app, "Ваши токены истекли. Производим вход в систему")
			if err := handleAuthProcess(app, session); err != nil {
				sendMessage(app, "Не удалось войти в систему. Используйте /start.")
				resetState(session)
				return
			} else {
				sendMessage(app, "Вход выполнен успешно!")
			}
		} else {
			sendMessage(app, "Ваш токен был успешно обновлен")
		}
	}
	next(app, session)
}
