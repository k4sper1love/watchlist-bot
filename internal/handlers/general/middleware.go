package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
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
			}
		} else {
			app.SendMessage("Ваш токен был успешно обновлен", nil)
		}
	}
	next(app, session)
}

func RequireAdmin(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if !session.IsAdmin {
		app.SendMessage("Недостаточный уровень прав", nil)
		session.ClearState()
		HandleMenuCommand(app, session)
		return
	}

	next(app, session)
}

func CheckBanned(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	isBanned, _ := postgres.IsUserBanned(session.TelegramID)

	if isBanned {
		app.SendMessage("❌ Вы заблокированы.\n Обратитесь к администратору для разблокировки.", nil)
		return
	}

	next(app, session)
}

func isAuth(session *models.Session) bool {
	if session.AccessToken == "" {
		return false
	}

	return true
}
