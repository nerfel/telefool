package middleware

import (
	"telefool/pkg/di"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PreventAddGroup(next Handler) Handler {
	return func(ctx *di.UpdateContext) {
		if ctx.Update.MyChatMember == nil {
			next(ctx)
			return
		}

		old := ctx.Update.MyChatMember.OldChatMember
		newm := ctx.Update.MyChatMember.NewChatMember

		if newm.User == nil || newm.User.ID != ctx.Bot.Self.ID {
			next(ctx)
			return
		}

		if old.Status == "left" && newm.Status == "member" {
			if ctx.Update.MyChatMember.From.UserName != ctx.Conf.AdminUserName {
				ctx.Bot.Request(tgbotapi.LeaveChatConfig{ChatID: ctx.Update.MyChatMember.Chat.ID})
			}

			return
		}

		if newm.Status == "administrator" {
			return
		}

		if newm.Status == "left" {
			return
		}

		next(ctx)
	}
}
