package stackable

type FuncHandlerWrapper[S ISharedState, L ILocalState] struct {
	Handler func(
		context *Context[S, L],
		next func() error,
	) error
}

func (h FuncHandlerWrapper[S, L]) Run(
	context *Context[S, L],
	next func() error,
) error {
	return h.Handler(
        context,
		next,
	)
}

func WrapFunc[S ISharedState, L ILocalState](handler func(
	context *Context[S, L],
	next func() error,
) error,
) FuncHandlerWrapper[S, L] {
	return FuncHandlerWrapper[S, L]{
		Handler: handler,
	}
}
