package middleware

import (
	"telefool/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsAdmin(next Handler) Handler {
	return func(update tgbotapi.Update, config *configs.Config, bot *tgbotapi.BotAPI) {
		if update.Message.From.UserName != config.AdminUserName {
			return
		}

		next(update, config, bot)
	}
}
