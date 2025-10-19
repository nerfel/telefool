package di

import (
	"telefool/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateContext struct {
	Update tgbotapi.Update
	Bot    *tgbotapi.BotAPI
	Conf   *configs.Config
}
