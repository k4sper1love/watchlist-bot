package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/api"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func StartHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := "Добро пожаловать в бота! Напишите команду /help чтобы узнать подробнее"
	user := models.User{TelegramID: update.Message.From.ID}
	postgres.GetDB().Create(&user)
	utils.SendMessage(bot, update, text)
}

func HelpHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := "Вот что умеет этот бот:\n" +
		"/start - приветственное сообщение\n" +
		"/help - узнать команды\n" +
		"/profile - получить профиль"
	utils.SendMessage(bot, update, text)
}

func ProfileHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var user models.User
	text := ""
	result := postgres.GetDB().First(&user, "telegram_id = ?", update.Message.From.ID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			text = "Профиль не найден"
		} else {
			text = "Произошла ошибка при получении профиля"
		}
	} else {
		text = fmt.Sprintf("Ваш профиль.\n Телеграмм ID: %d\nID записи в таблице: %d\nAccess Token: %s", user.TelegramID, user.ID, user.AccessToken)
	}
	utils.SendMessage(bot, update, text)
}

func LoginHandler(cfg *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	commandArgs := update.Message.CommandArguments()
	args := strings.Fields(commandArgs)

	text := ""

	if len(args) < 2 {
		text = "Неправильное количество аргументов"
	}

	var credentials struct {
		email    string
		password string
	}

	credentials.email = args[0]
	credentials.password = args[1]

	resp, err := api.SendRequest(cfg.BaseURL, "/auth/login", http.MethodPost, &credentials)
	if err != nil {
		text = "Ошибка при работе с API"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		text = "Успешный вход"
		postgres.GetDB().
			Model(&models.User{}).
			Where("telegram_id = ?", update.Message.From.ID).
			Update("access_token")

	} else {
		text = "Неудачный вход. Проверьте введенные данные"
	}

}
