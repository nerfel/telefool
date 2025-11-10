package middleware

import (
	"telefool/pkg/di"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsAdmin(next Handler) Handler {
	return func(ctx *di.UpdateContext, container *di.Container) {
		if ctx.Update.Message.From.UserName != ctx.Config.AdminUserName {
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "У вас не достаточно прав для данного функционала")
			msg.ReplyToMessageID = ctx.Update.Message.MessageID
			ctx.Bot.Send(msg)

			return
		}

		next(ctx, container)
	}
}
