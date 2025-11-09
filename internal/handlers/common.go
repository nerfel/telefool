package handlers

import (
	"telefool/configs"
	"telefool/internal/dialog"
	"telefool/internal/user"
	"telefool/pkg/di"
	"telefool/pkg/event"
	"telefool/pkg/memory"
	"telefool/pkg/middleware"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandlerDeps struct {
	Config        *configs.Config
	UserService   *user.UserService
	DialogService *dialog.DialogService
	EventBus      *event.Bus
	Bot           *tgbotapi.BotAPI
	Router        di.RouterInterface
	Memory        *memory.ShortTermMemory
}

type UpdateHandler struct {
	Config        *configs.Config
	UserService   *user.UserService
	DialogService *dialog.DialogService
	EventBus      *event.Bus
	Bot           *tgbotapi.BotAPI
	Router        di.RouterInterface
	Memory        *memory.ShortTermMemory
}

func NewUpdateHandler(deps *UpdateHandlerDeps) *UpdateHandler {
	return &UpdateHandler{
		Config:        deps.Config,
		UserService:   deps.UserService,
		DialogService: deps.DialogService,
		EventBus:      deps.EventBus,
		Bot:           deps.Bot,
		Router:        deps.Router,
		Memory:        deps.Memory,
	}
}

func (gmh *UpdateHandler) Handle() {
	gmh.Router.Register(CreateUserRoute, middleware.IsAdmin(CreateUserHandler))

	stack := middleware.Chain(
		middleware.PreventAddGroup,
		middleware.IgnoreEmpty,
		middleware.Logging,
	)
	handle := stack(func(ctx *di.UpdateContext) {
		gmh.Router.Serve(ctx)
	})

	for update := range gmh.Bot.ListenForWebhook("/") {
		handle(&di.UpdateContext{
			Update:   update,
			Bot:      gmh.Bot,
			EventBus: gmh.EventBus,
			Memory:   gmh.Memory,
			Config:   gmh.Config,
		})
	}
}
