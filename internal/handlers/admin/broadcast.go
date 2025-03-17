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

// HandleBroadcastCommand handles the command for initiating a broadcast message.
// Prompts the user to upload an image for the broadcast and sets the state to await the image.
func HandleBroadcastCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestBroadcastImage(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitBroadcastImage)
}

// HandleBroadcastProcess processes the broadcast workflow based on the current session state.
// Handles states like awaiting an image, text, pin confirmation, and final confirmation.
func HandleBroadcastProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitBroadcastImage:
		parser.ParseBroadcastImage(app, session, requestBroadcastMessage)

	case states.AwaitBroadcastText:
		parser.ParseBroadcastMessage(app, session, requestBroadcastMessage, requestBroadcastPin)

	case states.AwaitBroadcastPin:
		parser.ParseBroadcastPin(app, session, previewBroadcast)

	case states.AwaitBroadcastConfirm:
		parseBroadcastConfirm(app, session)
	}
}

// requestBroadcastMessage prompts the user to enter the text for the broadcast message.
// Sets the state to await the broadcast text.
func requestBroadcastMessage(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestBroadcastMessage(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitBroadcastText)
}

// requestBroadcastPin prompts the user to confirm whether the broadcast message should be pinned.
// Sets the state to await the pin confirmation.
func requestBroadcastPin(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestBroadcastPin(session), keyboards.SurveyAndCancel(session))
	session.SetState(states.AwaitBroadcastPin)
}

// previewBroadcast generates a preview of the broadcast message or image.
// If no content is provided, clears the states and returns to the menu.
func previewBroadcast(app models.App, session *models.Session) {
	if session.AdminState.Message == "" && session.AdminState.ImageURL == "" {
		app.SendMessage(messages.BroadcastEmpty(session), nil)
		clearStatesAndHandleMenu(app, session)
		return
	}

	if session.AdminState.ImageURL != "" {
		app.SendImage(session.AdminState.ImageURL, messages.BroadcastPreview(session), nil)
	} else {
		app.SendMessage(messages.BroadcastPreview(session), nil)
	}

	requestBroadcastConfirm(app, session)
}

// requestBroadcastConfirm prompts the user to confirm sending the broadcast.
// Includes the total number of recipients and sets the state to await confirmation.
func requestBroadcastConfirm(app models.App, session *models.Session) {
	count, err := postgres.GetUserCount(false)
	if err != nil {
		app.SendMessage(messages.RequestFailure(session), nil)
		clearStatesAndHandleMenu(app, session)
		return
	}

	app.SendMessage(messages.BroadcastConfirm(session, count), keyboards.BroadcastConfirm(session))
	session.SetState(states.AwaitBroadcastConfirm)
}

// parseBroadcastConfirm processes the final confirmation for sending the broadcast.
// Sends the broadcast message or image to all recipients and clears the session states.
func parseBroadcastConfirm(app models.App, session *models.Session) {
	if utils.ParseCallback(app.Update) != states.CallBroadcastSend {
		clearStatesAndHandleMenu(app, session)
		return
	}

	ids, err := postgres.GetTelegramIDs()
	if err != nil {
		app.SendMessage(messages.RequestFailure(session), nil)
		clearStatesAndHandleMenu(app, session)
		return
	}

	if session.AdminState.ImageURL != "" {
		app.SendBroadcastImage(ids, session.AdminState.NeedPin, session.AdminState.ImageURL, session.AdminState.Message, nil)
	} else {
		app.SendBroadcastMessage(ids, session.AdminState.NeedPin, session.AdminState.Message, nil)
	}

	clearStatesAndHandleMenu(app, session)
}

// clearStatesAndHandleMenu clears all session states and redirects the user to the main menu.
func clearStatesAndHandleMenu(app models.App, session *models.Session) {
	session.ClearAllStates()
	HandleMenuCommand(app, session)
}
