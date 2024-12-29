package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildStartMessage(app models.App, session *models.Session) string {
	part1 := translator.Translate(session.Lang, "welcomeMessageGreeting", map[string]interface{}{
		"Name": utils.ParseTelegramName(app.Upd),
	}, nil)

	part2 := translator.Translate(session.Lang, "welcomeMessageBody", map[string]interface{}{
		"Version": app.Vars.Version,
	}, nil)

	part3 := translator.Translate(session.Lang, "welcomeMessageCallToAction", nil, nil)

	msg := fmt.Sprintf("👋 <b>%s</b>\n\n%s 🚀\n\n%s", part1, part2, part3)

	return msg
}

func BuildHelpMessage(session *models.Session) string {
	msg := translator.Translate(session.Lang, "helpMessage", nil, nil)

	return msg
}

func BuildFeedbackMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "feedbackMessageHeader", nil, nil)
	part2 := translator.Translate(session.Lang, "feedbackMessageBody", nil, nil)
	part3 := translator.Translate(session.Lang, "feedbackCategoryChoice", nil, nil)

	msg := fmt.Sprintf("📝 <b>%s</b>\n\n%s😊\n\n%s", part1, part2, part3)

	return msg
}

func BuildMenuMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "mainMenu", nil, nil)
	part2 := translator.Translate(session.Lang, "choiceAction", nil, nil)

	msg := fmt.Sprintf("📋 <b>%s:</b>\n\n%s", part1, part2)

	return msg
}
