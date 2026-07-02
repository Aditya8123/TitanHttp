package router

import "github.com/Aditya8123/TitanHttp/internal/http"

// Handler defines the function signature for processing HTTP requests.
// Unlike the standard library's func(ResponseWriter, *Request), TitanHTTP
// uses a pure function approach that returns a constructed *http.Response.
type Handler func(req *http.Request) *http.Response

// Router manages the registration and dispatching of HTTP routes.
type Router struct {
	// routes maps an exact path string to a Handler.
	// This will be evolved into a Radix tree in future subtasks for parameters/wildcards.
	routes map[string]Handler
}

// NewRouter initializes and returns a new Router instance.
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
	}
}

// AddRoute registers a new handler for the given exact path.
func (r *Router) AddRoute(path string, handler Handler) {
	r.routes[path] = handler
}

// ServeHTTP attempts to find a matching route and execute its handler.
// If no route matches, it returns a generic 404 Not Found response.
func (r *Router) ServeHTTP(req *http.Request) *http.Response {
	// Basic exact matching for now.
	if handler, exists := r.routes[req.Path]; exists {
		return handler(req)
	}

	// Fallback to 404
	return http.NewResponse404()
}
