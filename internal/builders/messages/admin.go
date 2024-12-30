package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildAdminUserListMessage(session *models.Session, users []models.Session) string {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize
	lastPage := session.AdminState.LastPage
	totalRecords := session.AdminState.TotalRecords

	foundMsg := translator.Translate(session.Lang, "foundUsers", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b> %d\n\n", foundMsg, totalRecords)

	for i, user := range users {
		itemID := utils.GetItemID(i, currentPage, pageSize)

		numberEmoji := numberToEmoji(itemID)

		roleMsg := translator.Translate(session.Lang, user.Role.String(), nil, nil)
		msg += fmt.Sprintf("<b>%s %s</b>\n", numberEmoji, roleMsg)

		msg += fmt.Sprintf("Telegram ID: <code>%d</code>\n", user.TelegramID)

		if user.TelegramUsername != "" {
			msg += fmt.Sprintf("Username: @%s\n", user.TelegramUsername)
		}

		bannedMsg := translator.Translate(session.Lang, "banned", nil, nil)
		msg += fmt.Sprintf("%s: %s\n", bannedMsg, boolToEmoji(user.IsBanned))

		msg += "\n"
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": currentPage,
		"LastPage":    lastPage,
	}, nil)

	msg += fmt.Sprintf("<b>📄 %s</b>\n", pageMsg)
	msg += translator.Translate(session.Lang, "choiceUserForDetails", nil, nil)

	return msg
}

func BuildAdminUserDetailMessage(session *models.Session, user *models.Session) string {
	detailsMsg := translator.Translate(session.Lang, "sessionDetails", nil, nil)
	msg := fmt.Sprintf("💻 <b>%s:</b>\n", detailsMsg)

	msg += fmt.Sprintf("<b>Telegram ID:</b> <code>%d</code>\n", user.TelegramID)

	if user.TelegramUsername != "" {
		msg += fmt.Sprintf("<b>Username:</b> @%s\n", user.TelegramUsername)
	}

	roleMsg := translator.Translate(session.Lang, "role", nil, nil)
	roleValueMsg := translator.Translate(session.Lang, user.Role.String(), nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", roleMsg, roleValueMsg)

	bannedMsg := translator.Translate(session.Lang, "banned", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", bannedMsg, boolToEmoji(user.IsBanned))

	languageMsg := translator.Translate(session.Lang, "language", nil, nil)
	language := ""
	if user.Lang == "" {
		language = translator.Translate(session.Lang, "empty", nil, nil)
	} else {
		language = user.Lang
	}
	msg += fmt.Sprintf("<b>%s:</b> %s\n", languageMsg, language)

	stateMsg := translator.Translate(session.Lang, "state", nil, nil)
	state := ""
	if user.State == "" {
		state = translator.Translate(session.Lang, "empty", nil, nil)
	} else {
		state = translator.Translate(session.Lang, "empty", nil, nil)
	}
	msg += fmt.Sprintf("<b>%s:</b> %s\n", stateMsg, state)

	createdMsg := translator.Translate(session.Lang, "created", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", createdMsg, user.CreatedAt.Format("02.01.2006 15:04"))

	updatedMsg := translator.Translate(session.Lang, "updated", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n\n", updatedMsg, user.UpdatedAt.Format("02.01.2006 15:04"))

	detailsMsg = translator.Translate(session.Lang, "apiDetails", nil, nil)
	msg += fmt.Sprintf("🌐 <b>%s:</b>\n", detailsMsg)

	msg += fmt.Sprintf("<b>ID:</b> <code>%d</code>\n", user.User.ID)

	nameMsg := translator.Translate(session.Lang, "name", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> <code>%s</code>\n", nameMsg, user.User.Username)

	if user.User.Email != "" {
		msg += fmt.Sprintf("<b>Email:</b> <code>%s</code>\n", user.User.Email)
	}

	createdMsg = translator.Translate(session.Lang, "created", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", createdMsg, user.User.CreatedAt.Format("02.01.2006 15:04"))

	return msg
}

func BuildFeedbackListMessage(session *models.Session, feedbacks []models.Feedback) string {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize
	lastPage := session.AdminState.LastPage
	totalRecords := session.AdminState.TotalRecords

	foundMsg := translator.Translate(session.Lang, "foundFeedbacks", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b> %d\n\n", foundMsg, totalRecords)

	for i, feedback := range feedbacks {
		itemID := utils.GetItemID(i, currentPage, pageSize)

		numberEmoji := numberToEmoji(itemID)

		categoryMsg := translator.Translate(session.Lang, feedback.Category, nil, nil)
		msg += fmt.Sprintf("<b>%s %s</b>\n", numberEmoji, categoryMsg)

		idMsg := translator.Translate(session.Lang, "user", nil, nil)
		msg += fmt.Sprintf("%s: <code>%d</code>\n", idMsg, feedback.TelegramID)

		createdMsg := translator.Translate(session.Lang, "created", nil, nil)
		msg += fmt.Sprintf("%s: %s\n\n", createdMsg, feedback.CreatedAt.Format("02.01.2006 15:04"))
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": currentPage,
		"LastPage":    lastPage,
	}, nil)

	msg += fmt.Sprintf("<b>📄 %s</b>\n", pageMsg)
	msg += translator.Translate(session.Lang, "choiceUserForDetails", nil, nil)

	return msg
}

func BuildFeedbackDetailMessage(session *models.Session, feedback *models.Feedback) string {
	detailsMsg := translator.Translate(session.Lang, "feedbackDetails", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b>\n\n", detailsMsg)

	idMsg := translator.Translate(session.Lang, "user", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> <code>%d</code>\n", idMsg, feedback.TelegramID)

	categoryMsg := translator.Translate(session.Lang, "category", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", categoryMsg, feedback.Category)

	feedbackMsg := translator.Translate(session.Lang, "message", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b>\n<pre>%s</pre>\n", feedbackMsg, feedback.Message)

	createdMsg := translator.Translate(session.Lang, "created", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", createdMsg, feedback.CreatedAt.Format("02.01.2006 15:04"))

	return msg
}
