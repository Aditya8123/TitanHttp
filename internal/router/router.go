package router

import "github.com/Aditya8123/TitanHttp/internal/http"

// Handler defines the function signature for processing HTTP requests.
// Unlike the standard library's func(ResponseWriter, *Request), TitanHTTP
// uses a pure function approach that returns a constructed *http.Response.
type Handler func(req *http.Request) *http.Response

// Router manages the registration and dispatching of HTTP routes.
type Router struct {
	// routes maps an HTTP Method to a map of exact path strings to Handlers.
	// This will be evolved into a Radix tree in future subtasks for parameters/wildcards.
	routes map[http.Method]map[string]Handler
}

// NewRouter initializes and returns a new Router instance.
func NewRouter() *Router {
	return &Router{
		routes: make(map[http.Method]map[string]Handler),
	}
}

// AddRoute registers a new handler for the given HTTP method and exact path.
func (r *Router) AddRoute(method http.Method, path string, handler Handler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]Handler)
	}
	r.routes[method][path] = handler
}

// Get is a convenience method for registering a GET route.
func (r *Router) Get(path string, handler Handler) {
	r.AddRoute(http.MethodGet, path, handler)
}

// Post is a convenience method for registering a POST route.
func (r *Router) Post(path string, handler Handler) {
	r.AddRoute(http.MethodPost, path, handler)
}

// ServeHTTP attempts to find a matching route and execute its handler.
// If the path matches but the method does not, it returns a 405 Method Not Allowed.
// If no route matches the path at all, it returns a 404 Not Found.
func (r *Router) ServeHTTP(req *http.Request) *http.Response {
	methodRoutes, methodExists := r.routes[req.Method]
	if methodExists {
		if handler, exists := methodRoutes[req.Path]; exists {
			return handler(req)
		}
	}

	// 405 check: Does this path exist under ANY other method?
	for _, mRoutes := range r.routes {
		if _, exists := mRoutes[req.Path]; exists {
			return http.NewResponse405()
		}
	}

	// Fallback to 404
	return http.NewResponse404()
}
