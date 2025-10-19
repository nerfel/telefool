package middleware

import (
	"telefool/pkg/di"
)

func IgnoreEmpty(next Handler) Handler {
	return func(ctx *di.UpdateContext) {
		if ctx.Update.Message.Text == "" {
			return
		}

		next(ctx)
	}
}
