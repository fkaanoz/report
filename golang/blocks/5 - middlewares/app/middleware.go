package app

type Middleware func(handler Handler) Handler

// Wrap will wrap the given handler with middlewares.
func Wrap(handler Handler, mv ...Middleware) Handler {

	for _, middleware := range mv {
		if middleware != nil {
			handler = middleware(handler)
		}
	}

	return handler
}
