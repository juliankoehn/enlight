package enlight

// MiddlewareFunc defines a function to process middleware
type MiddlewareFunc func(HandleFunc) HandleFunc

func applyMiddleware(h HandleFunc, middleware ...MiddlewareFunc) HandleFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
