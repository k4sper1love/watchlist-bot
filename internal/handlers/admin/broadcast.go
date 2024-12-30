package admin

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleBroadcastCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "requestBroadcastImage", nil, nil)

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAdminBroadcastAwaitingImage)
}

func HandleBroadcastProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAdminBroadcastAwaitingImage:
		parseBroadcastImage(app, session)

	case states.ProcessAdminBroadcastAwaitingText:
		parseBroadcastMessage(app, session)

	case states.ProcessAdminBroadcastAwaitingConfirm:
		parseBroadcastConfirm(app, session)
	}
}

func parseBroadcastImage(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		requestBroadcastMessage(app, session)
		return
	}

	image, err := utils.ParseImageFromMessage(app.Bot, app.Upd)
	if err != nil {
		handleBroadcastImageError(app, session)
		return
	}

	imageURL, err := watchlist.UploadImage(app, image)
	if err != nil {
		handleBroadcastImageError(app, session)
		return
	}

	session.AdminState.FeedbackImageURL = imageURL

	requestBroadcastMessage(app, session)
}

func requestBroadcastMessage(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "requestBroadcastMessage", nil, nil)

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAdminBroadcastAwaitingText)
}

func parseBroadcastMessage(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		requestBroadcastConfirm(app, session)
		return
	}

	msg := utils.ParseMessageString(app.Upd)

	session.AdminState.FeedbackMessage = msg

	if session.AdminState.FeedbackImageURL != "" {
		app.SendImage(session.AdminState.FeedbackImageURL, msg, nil)
	} else {
		app.SendMessage(msg, nil)
	}

	requestBroadcastConfirm(app, session)
}

func requestBroadcastConfirm(app models.App, session *models.Session) {
	if session.AdminState.FeedbackMessage == "" && session.AdminState.FeedbackImageURL == "" {
		emptyMsg := translator.Translate(session.Lang, "broadcastEmpty", nil, nil)
		app.SendMessage(emptyMsg, nil)
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	count, err := postgres.GetUserCounts()
	if err != nil {
		msg := translator.Translate(session.Lang, "requestFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	countMsg := translator.Translate(session.Lang, "recipientCount", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>: %d", countMsg, count)

	keyboard := keyboards.BuildBroadcastConfirmKeyboard(session)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAdminBroadcastAwaitingConfirm)
}

func parseBroadcastConfirm(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackAdminBroadcastSend:
		telegramIDs, err := postgres.GetAllTelegramID()
		if err != nil {
			msg := translator.Translate(session.Lang, "requestFailure", nil, nil)
			app.SendMessage(msg, nil)
			session.ClearAllStates()
			HandleMenuCommand(app, session)
			return
		}

		if session.AdminState.FeedbackImageURL != "" {
			app.SendBroadcastImage(telegramIDs, session.AdminState.FeedbackImageURL, session.AdminState.FeedbackMessage, nil)
		} else {
			app.SendBroadcastMessage(telegramIDs, session.AdminState.FeedbackMessage, nil)
		}
	}

	session.ClearAllStates()
	HandleMenuCommand(app, session)
}

func handleBroadcastImageError(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "getImageFailure", nil, nil)
	app.SendMessage(msg, nil)
	requestBroadcastMessage(app, session)
}
