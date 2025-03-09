package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"log/slog"
)

func Auth(app models.App, session *models.Session) bool {
	if IsBanned(app, session) {
		return false
	}

	if isAuthenticated(app, session) || attemptLoginOrRegister(app, session) == nil {
		return true
	}

	app.SendMessage(messages.AuthFailure(session), nil)
	session.ClearAllStates()
	return false
}

func RequireAuth(app models.App, session *models.Session, next func(app models.App, session *models.Session)) {
	if Auth(app, session) {
		next(app, session)
	}
}

func RequireRole(app models.App, session *models.Session, next func(models.App, *models.Session), role roles.Role) {
	if session.Role.HasAccess(role) {
		next(app, session)
		return
	}

	app.SendMessage(messages.PermissionsNotEnough(session), keyboards.BuildKeyboardWithBack(session, ""))
	session.ClearState()
}

func IsBanned(app models.App, session *models.Session) bool {
	if !postgres.IsUserBanned(session.TelegramID) {
		return false
	}

	app.SendMessage(messages.Banned(session), nil)
	return true
}

func isAuthenticated(app models.App, session *models.Session) bool {
	return session.AccessToken != "" &&
		(watchlist.IsTokenValid(app, session.AccessToken) ||
			(session.RefreshToken != "" && watchlist.RefreshAccessToken(app, session) == nil))
}

func attemptLoginOrRegister(app models.App, session *models.Session) error {
	if err := watchlist.Login(app, session); err == nil {
		return nil
	}

	if err := watchlist.Register(app, session); err != nil {
		sl.Log.Error("failed to login/register", slog.Any("error", err), slog.Int("telegram_id", session.TelegramID))
		return err
	}

	app.SendMessage(messages.RegistrationSuccess(session), nil)
	return nil
}
