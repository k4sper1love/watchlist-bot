package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
)

func requireAuth(app config.App, session *models.Session, next func(config.App, *models.Session)) {
	if !session.IsLogged {
		utils.SendMessage(app.Bot, app.Upd, "Для выполнения этой команды нужно быть авторизованным. Используйте /register или /login")
		return
	}

	if !watchlist.IsTokenValid(app, session.AccessToken) {
		if err := watchlist.RefreshAccessToken(app, session); err != nil {
			utils.SendMessage(app.Bot, app.Upd, "Ваши токены истекли. Авторизуйтесь заново с помощью /login")
			session.IsLogged = false
			return
		} else {
			utils.SendMessage(app.Bot, app.Upd, "Ваш токен был успешно обновлен")
		}
	}

	next(app, session)
}

func requireNoAuth(app config.App, session *models.Session, next func(config.App, *models.Session)) {
	if session.IsLogged {
		utils.SendMessage(app.Bot, app.Upd, "Вы уже авторизованы")
		return
	}

	next(app, session)
}
