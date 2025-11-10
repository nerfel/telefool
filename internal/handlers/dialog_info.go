package handlers

import (
	"fmt"
	"telefool/pkg/di"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const dialogInfoCommand = "get-dialog-info"

func DialogInfoRoute(ctx *di.UpdateContext) bool {
	return ctx.Update.Message.Text == dialogInfoCommand
}

func DialogInfoHandler(ctx *di.UpdateContext, container *di.Container) {
	ctx.Bot.Send(
		tgbotapi.NewMessage(
			ctx.Update.Message.Chat.ID,
			fmt.Sprintf(
				"Chat Title: %s,\nChat ID:%d",
				ctx.Update.Message.Chat.Title,
				ctx.Update.Message.Chat.ID,
			),
		),
	)
}
