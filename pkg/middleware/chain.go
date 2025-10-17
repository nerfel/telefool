package middleware

import (
	"telefool/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(update tgbotapi.Update, config *configs.Config, bot *tgbotapi.BotAPI)
type Middleware func(Handler) Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(h Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}
