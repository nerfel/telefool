package middleware

import (
	"telefool/pkg/di"
	"telefool/pkg/event"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PreventAddGroup(next Handler) Handler {
	return func(ctx *di.UpdateContext, container *di.Container) {
		if ctx.Update.MyChatMember == nil {
			next(ctx, container)
			return
		}

		old := ctx.Update.MyChatMember.OldChatMember
		newm := ctx.Update.MyChatMember.NewChatMember

		if newm.User == nil || newm.User.ID != ctx.Bot.Self.ID {
			next(ctx, container)
			return
		}

		if old.Status == "left" && newm.Status == "member" {
			if ctx.Update.MyChatMember.From.UserName != ctx.Config.AdminUserName {
				ctx.Bot.Request(tgbotapi.LeaveChatConfig{ChatID: ctx.Update.MyChatMember.Chat.ID})
			}

			ctx.EventBus.Publish(event.Event{Type: event.EventAddToGroup, Data: ctx.Update}) // side effect

			return
		}

		if newm.Status == "administrator" {
			return
		}

		if newm.Status == "left" {
			ctx.EventBus.Publish(event.Event{Type: event.EventRemoveFromGroup, Data: ctx.Update}) // side effect
			return
		}

		next(ctx, container)
	}
}
