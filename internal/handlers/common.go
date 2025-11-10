package handlers

import (
	"telefool/configs"
	"telefool/pkg/di"
	"telefool/pkg/event"
	"telefool/pkg/memory"
	"telefool/pkg/middleware"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandlerDeps struct {
	Config    *configs.Config
	EventBus  *event.Bus
	Bot       *tgbotapi.BotAPI
	Router    di.RouterInterface
	Container *di.Container
	Memory    *memory.ShortTermMemory
}

type UpdateHandler struct {
	Config    *configs.Config
	EventBus  *event.Bus
	Bot       *tgbotapi.BotAPI
	Router    di.RouterInterface
	Container *di.Container
	Memory    *memory.ShortTermMemory
}

func NewUpdateHandler(deps *UpdateHandlerDeps) *UpdateHandler {
	return &UpdateHandler{
		Config:    deps.Config,
		EventBus:  deps.EventBus,
		Bot:       deps.Bot,
		Router:    deps.Router,
		Container: deps.Container,
		Memory:    deps.Memory,
	}
}

func (gmh *UpdateHandler) Handle() {
	gmh.Router.Register(CreateUserRoute, middleware.IsAdmin(CreateUserHandler))

	stack := middleware.Chain(
		middleware.PreventAddGroup,
		middleware.IgnoreEmpty,
		middleware.Logging,
	)
	handle := stack(func(ctx *di.UpdateContext, container *di.Container) {
		gmh.Router.Serve(ctx, container)
	})

	for update := range gmh.Bot.ListenForWebhook("/") {
		handle(&di.UpdateContext{
			Update:   update,
			Bot:      gmh.Bot,
			EventBus: gmh.EventBus,
			Memory:   gmh.Memory,
			Config:   gmh.Config,
		}, gmh.Container)
	}
}
