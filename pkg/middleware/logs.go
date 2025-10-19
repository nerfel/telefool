package middleware

import (
	"log"
	"telefool/pkg/di"
)

func Logging(next Handler) Handler {
	return func(ctx *di.UpdateContext) {
		log.Printf("message_from: [%s], message_text: %s", ctx.Update.Message.From, ctx.Update.Message.Text)
		next(ctx)
	}
}
