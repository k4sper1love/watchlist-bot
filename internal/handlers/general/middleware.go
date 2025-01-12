package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func RequireAuth(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if !isAuth(session) {
		if err := HandleAuthProcess(app, session); err != nil {
			msg := translator.Translate(session.Lang, "authRequest", nil, nil)
			app.SendMessage(msg, nil)
			session.ClearState()
			return
		}
	} else if !watchlist.IsTokenValid(app, session.AccessToken) {
		if err := watchlist.RefreshAccessToken(app, session); err != nil {
			//msg := translator.Translate(session.Lang, "authExpired", nil, nil)
			//app.SendMessage(msg, nil)
			if err := HandleAuthProcess(app, session); err != nil {
				msg := translator.Translate(session.Lang, "authFailure", nil, nil)
				app.SendMessage(msg, nil)
				session.ClearState()
				return
			}

		} else {
			//msg := translator.Translate(session.Lang, "authUpdated", nil, nil)
			//app.SendMessage(msg, nil)
		}
	}
	next(app, session)
}

func RequireRole(app models.App, session *models.Session, next func(models.App, *models.Session), role roles.Role) {
	if !session.Role.HasAccess(role) {
		msg := translator.Translate(session.Lang, "permissionsNotEnough", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearState()
		HandleMenuCommand(app, session)
		return
	}

	RequireAuth(app, session, next)
}

func CheckBanned(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	isBanned, _ := postgres.IsUserBanned(session.TelegramID)

	if isBanned {
		part1 := translator.Translate(session.Lang, "bannedHeader", nil, nil)
		part2 := translator.Translate(session.Lang, "bannedBody", nil, nil)

		msg := fmt.Sprintf("❌ %s\n\n%s", part1, part2)

		app.SendMessage(msg, nil)
		return
	}

	RequireAuth(app, session, next)
}

func isAuth(session *models.Session) bool {
	if session.AccessToken == "" {
		return false
	}

	return true
}
