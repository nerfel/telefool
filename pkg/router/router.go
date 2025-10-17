package router

import (
	"telefool/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateContext struct {
	Update tgbotapi.Update
	Bot    *tgbotapi.BotAPI
}

type Route struct {
	Match  func(ctx *UpdateContext, conf configs.Config) bool
	Handle func(ctx *UpdateContext)
}

type Router struct {
	config configs.Config
	bot    *tgbotapi.BotAPI
	routes []Route
}

func NewUpdateRouter(config *configs.Config, bot *tgbotapi.BotAPI) *Router {
	return &Router{config: *config, bot: bot}
}

func (r *Router) Register(match func(*UpdateContext, configs.Config) bool, handle func(ctx *UpdateContext)) {
	r.routes = append(r.routes, Route{match, handle})
}

func (r *Router) Handle(update tgbotapi.Update) {
	ctx := &UpdateContext{Update: update, Bot: r.bot}
	for _, route := range r.routes {
		if route.Match(ctx, r.config) {
			route.Handle(ctx)
			return
		}
	}
}
