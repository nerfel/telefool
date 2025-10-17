package handlers

import (
	"telefool/configs"
	"telefool/internal/user"
	"telefool/pkg/middleware"
	"telefool/pkg/router"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandlerDeps struct {
	Config      *configs.Config
	UserService *user.UserService
	Bot         *tgbotapi.BotAPI
}

type UpdateHandler struct {
	Config      *configs.Config
	UserService *user.UserService
	Bot         *tgbotapi.BotAPI
}

func NewUpdateHandler(deps *UpdateHandlerDeps) *UpdateHandler {
	return &UpdateHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
		Bot:         deps.Bot,
	}
}

func (mh *UpdateHandler) Handle() {
	r := router.NewUpdateRouter(mh.Config, mh.Bot)

	stack := middleware.Chain(
		middleware.PreventAddGroup,
		middleware.IgnoreEmpty,
		middleware.Logging,
	)
	handle := stack(func(update tgbotapi.Update, config *configs.Config, bot *tgbotapi.BotAPI) {
		r.Handle(update)
	})

	for update := range mh.Bot.ListenForWebhook("/") {
		handle(
			update,
			mh.Config,
			mh.Bot,
		)
	}
}
