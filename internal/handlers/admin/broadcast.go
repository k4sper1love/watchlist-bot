package admin

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleBroadcastCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildRequestBroadcastImageMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessAdminBroadcastAwaitingImage)
}

func HandleBroadcastProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAdminBroadcastAwaitingImage:
		parser.ParseBroadcastImage(app, session, requestBroadcastMessage)

	case states.ProcessAdminBroadcastAwaitingText:
		parser.ParseBroadcastMessage(app, session, requestBroadcastMessage, requestBroadcastPin)

	case states.ProcessAdminBroadcastAwaitingPin:
		parser.ParseBroadcastPin(app, session, previewBroadcast)

	case states.ProcessAdminBroadcastAwaitingConfirm:
		parseBroadcastConfirm(app, session)
	}
}

func requestBroadcastMessage(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildRequestBroadcastMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessAdminBroadcastAwaitingText)
}

func requestBroadcastPin(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildRequestBroadcastPinMessage(session), keyboards.BuildKeyboardWithSurveyAndCancel(session))
	session.SetState(states.ProcessAdminBroadcastAwaitingPin)
}

func previewBroadcast(app models.App, session *models.Session) {
	if session.AdminState.ImageURL != "" {
		app.SendImage(session.AdminState.ImageURL, messages.BuildBroadcastPreviewMessage(session), nil)
	} else {
		app.SendMessage(messages.BuildBroadcastPreviewMessage(session), nil)
	}

	requestBroadcastConfirm(app, session)
}

func requestBroadcastConfirm(app models.App, session *models.Session) {
	if session.AdminState.Message == "" && session.AdminState.ImageURL == "" {
		app.SendMessage(messages.BuildBroadcastEmptyMessage(session), nil)
		clearStatesAndHandleMenu(app, session)
		return
	}

	count, err := postgres.GetUserCount(false)
	if err != nil {
		app.SendMessage(messages.BuildRequestFailureMessage(session), nil)
		clearStatesAndHandleMenu(app, session)
		return
	}

	app.SendMessage(messages.BuildBroadcastConfirmMessage(session, count), keyboards.BuildBroadcastConfirmKeyboard(session))
	session.SetState(states.ProcessAdminBroadcastAwaitingConfirm)
}

func parseBroadcastConfirm(app models.App, session *models.Session) {
	if utils.ParseCallback(app.Update) != states.CallbackAdminBroadcastSend {
		clearStatesAndHandleMenu(app, session)
		return
	}

	ids, err := postgres.GetTelegramIDs()
	if err != nil {
		app.SendMessage(messages.BuildRequestFailureMessage(session), nil)
		clearStatesAndHandleMenu(app, session)
		return
	}

	if session.AdminState.ImageURL != "" {
		app.SendBroadcastImage(ids, session.AdminState.NeedFeedbackPin, session.AdminState.ImageURL, session.AdminState.Message, nil)
	} else {
		app.SendBroadcastMessage(ids, session.AdminState.NeedFeedbackPin, session.AdminState.Message, nil)
	}

	clearStatesAndHandleMenu(app, session)
}

func clearStatesAndHandleMenu(app models.App, session *models.Session) {
	session.ClearAllStates()
	HandleMenuCommand(app, session)
}
