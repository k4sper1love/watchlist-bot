package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func Auth(app models.App, session *models.Session) bool {
	if IsBanned(app, session) {
		return false
	}

	if !isAuth(session) {
		if err := HandleAuthProcess(app, session); err != nil {
			msg := translator.Translate(session.Lang, "authFailure", nil, nil)
			app.SendMessage(msg, nil)
			session.ClearAllStates()
			return false
		}
	} else if !watchlist.IsTokenValid(app, session.AccessToken) {
		if err := watchlist.RefreshAccessToken(app, session); err != nil {
			session.AccessToken = ""
			if err = HandleAuthProcess(app, session); err != nil {
				msg := translator.Translate(session.Lang, "authFailure", nil, nil)
				app.SendMessage(msg, nil)
				session.ClearAllStates()
				return false
			}

		}
	}

	return true
}

func RequireAuth(app models.App, session *models.Session, next func(app models.App, session *models.Session)) {
	if ok := Auth(app, session); ok {
		next(app, session)
	}
}

func RequireRole(app models.App, session *models.Session, next func(models.App, *models.Session), role roles.Role) {
	if !session.Role.HasAccess(role) {
		msg := translator.Translate(session.Lang, "permissionsNotEnough", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearState()
		HandleMenuCommand(app, session)
		return
	}

	next(app, session)
}

func IsBanned(app models.App, session *models.Session) bool {
	isBanned := postgres.IsUserBanned(session.TelegramID)

	if isBanned {
		part1 := translator.Translate(session.Lang, "bannedHeader", nil, nil)
		part2 := translator.Translate(session.Lang, "bannedBody", nil, nil)

		msg := fmt.Sprintf("❌ %s\n\n%s", part1, part2)

		app.SendMessage(msg, nil)
		return true
	}

	return false
}

func isAuth(session *models.Session) bool {
	if session.AccessToken == "" {
		return false
	}

	return true
}
