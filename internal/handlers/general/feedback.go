package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleFeedbackCommand(app models.App, session *models.Session) {
	msg := messages.BuildFeedbackMessage()

	keyboard := keyboards.BuildFeedbackKeyboard()

	app.SendMessage(msg, keyboard)
}

func HandleFeedbackButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	var category string
	switch callback {
	case states.CallbackFeedbackCategorySuggestions:
		category = "Предложения"
	case states.CallbackFeedbackCategoryBugs:
		category = "Ошибки"
	case states.CallbackFeedbackCategoryOther:
		category = "Другие вопросы"
	}

	session.FeedbackState.Category = category

	msg := fmt.Sprintf("📄 <b>Вы выбрали категорию:</b> %s\n\nПожалуйста, напишите ваш фидбек.", category)

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessFeedbackAwaitingMessage)
}

func HandleFeedbackProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		HandleFeedbackCommand(app, session)
	}

	switch session.State {
	case states.ProcessFeedbackAwaitingMessage:
		parseFeedbackMessage(app, session)
	}
}

func parseFeedbackMessage(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddBack("").Build()

	session.FeedbackState.Message = utils.ParseMessageString(app.Upd)

	err := postgres.SaveFeedbackToDatabase(session.TelegramID, session.FeedbackState.Category, session.FeedbackState.Message)
	if err != nil {
		app.SendMessage("❌ Ошибка сохранения фидбека. Попробуйте позже.", keyboard)
		return
	}

	app.SendMessage("✅ Спасибо за ваш фидбек! Мы ценим ваше мнение.", keyboard)

	session.ClearAllStates()
}
