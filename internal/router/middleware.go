package router

// Middleware is a function that wraps a Handler, returning a new Handler.
// This allows logic to execute before or after the main Handler, or even short-circuit it.
type Middleware func(Handler) Handler

// Chain combines multiple middlewares into a single Middleware.
// The middlewares are applied in the order they are passed, meaning the first
// middleware in the list is the outermost one (executes first on the request path,
// and last on the response path).
func Chain(middlewares ...Middleware) Middleware {
	return func(final Handler) Handler {
		// Loop backwards so that the first middleware ends up being the outermost wrapper
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
