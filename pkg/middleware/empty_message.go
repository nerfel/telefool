package middleware

import (
	"telefool/pkg/di"
)

func IgnoreEmpty(next Handler) Handler {
	return func(ctx *di.UpdateContext, container *di.Container) {
		if ctx.Update.Message == nil || ctx.Update.Message.Text == "" {
			return
		}

		next(ctx, container)
	}
}
