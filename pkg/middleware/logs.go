package middleware

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Logging(next Handler) Handler {
	return func(update tgbotapi.Update) {
		log.Printf("message_from: [%s], message_text: %s", update.Message.From, update.Message.Text)
		next(update)
	}
}
