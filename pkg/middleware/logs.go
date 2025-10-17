package middleware

import (
	"log"
	"telefool/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Logging(next Handler) Handler {
	return func(update tgbotapi.Update, config *configs.Config, bot *tgbotapi.BotAPI) {
		log.Printf("message_from: [%s], message_text: %s", update.Message.From, update.Message.Text)
		next(update, config, bot)
	}
}
