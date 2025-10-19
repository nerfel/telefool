package middleware

import (
	"telefool/pkg/di"
)

func IsAdmin(next Handler) Handler {
	return func(ctx *di.UpdateContext) {
		if ctx.Update.Message.From.UserName != ctx.Conf.AdminUserName {
			return
		}

		next(ctx)
	}
}
