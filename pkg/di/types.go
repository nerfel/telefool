package di

import (
	"telefool/configs"
	"telefool/pkg/event"
	"telefool/pkg/memory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateContext struct {
	Update       tgbotapi.Update
	Bot          *tgbotapi.BotAPI
	Config       *configs.Config
	EventBus     *event.Bus
	Memory       *memory.ShortTermMemory
	RoutePayload any
}
