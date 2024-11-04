package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
)

func RequireAuth(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if !isAuth(session) {
		app.SendMessage("Для выполнения этой команды нужно быть авторизованным. Используйте /start", nil)
		session.ClearState()
		return
	}

	if !watchlist.IsTokenValid(app, session.AccessToken) {
		if err := watchlist.RefreshAccessToken(app, session); err != nil {
			app.SendMessage("Ваши токены истекли. Производим вход в систему", nil)
			if err := HandleAuthProcess(app, session); err != nil {
				app.SendMessage("Не удалось войти в систему. Используйте /start.", nil)
				session.ClearState()
				return
			} else {
				app.SendMessage("Вход выполнен успешно!", nil)
			}
		} else {
			app.SendMessage("Ваш токен был успешно обновлен", nil)
		}
	}
	next(app, session)
}

func isAuth(session *models.Session) bool {
	if session.AccessToken == "" {
		return false
	}

	return true
}
