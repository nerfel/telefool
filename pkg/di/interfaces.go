package di

type RouterInterface interface {
	Register(func(*UpdateContext) bool, func(*UpdateContext, *Container))
	Serve(*UpdateContext, *Container)
}
