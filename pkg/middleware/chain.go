package middleware

import (
	"telefool/pkg/di"
)

type Handler func(ctx *di.UpdateContext)
type Middleware func(Handler) Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(h Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}
