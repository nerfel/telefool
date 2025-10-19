package handlers

import (
	"telefool/configs"
	"telefool/internal/user"
	"telefool/pkg/di"
	"telefool/pkg/middleware"
	"telefool/pkg/router"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandlerDeps struct {
	Config      *configs.Config
	UserService *user.UserService
	Bot         *tgbotapi.BotAPI
	Router      *router.Router
}

type UpdateHandler struct {
	Config      *configs.Config
	UserService *user.UserService
	Bot         *tgbotapi.BotAPI
	Router      *router.Router
}

func NewUpdateHandler(deps *UpdateHandlerDeps) *UpdateHandler {
	return &UpdateHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
		Bot:         deps.Bot,
		Router:      deps.Router,
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
			Update: update,
			Bot:    gmh.Bot,
			Conf:   gmh.Config,
		})
	}
}
