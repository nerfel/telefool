package router

import (
	"telefool/internal/handlers"
	"telefool/pkg/di"
)

type Route struct {
	Match  func(ctx *di.UpdateContext) bool
	Handle func(ctx *di.UpdateContext)
}

type Router struct {
	routes []Route
}

func NewUpdateRouter() *Router {
	return &Router{}
}

func (r *Router) Register(match func(*di.UpdateContext) bool, handle func(ctx *di.UpdateContext)) {
	r.routes = append(r.routes, Route{match, handle})
}

func (r *Router) Serve(ctx *di.UpdateContext) {
	for _, route := range r.routes {
		if route.Match(ctx) {
			route.Handle(ctx)
			return
		}
	}

	handlers.FallBackGPTHandle(ctx)
}
