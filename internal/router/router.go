package router

import "github.com/Aditya8123/TitanHttp/internal/http"

// Handler defines the function signature for processing HTTP requests.
type Handler func(req *http.Request) *http.Response

// Router manages the registration and dispatching of HTTP routes using Radix trees.
type Router struct {
	// trees maintains a Radix tree (prefix tree) for each HTTP method.
	trees map[http.Method]*node
}

// NewRouter initializes and returns a new Router instance.
func NewRouter() *Router {
	return &Router{
		trees: make(map[http.Method]*node),
	}
}

// AddRoute registers a new handler for the given HTTP method and pattern.
func (r *Router) AddRoute(method http.Method, pattern string, handler Handler) {
	if r.trees[method] == nil {
		r.trees[method] = &node{}
	}
	r.trees[method].insert(pattern, handler)
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
// Extracts dynamic path parameters and injects them into the Request.
func (r *Router) ServeHTTP(req *http.Request) *http.Response {
	tree, methodExists := r.trees[req.Method]
	if methodExists {
		handler, params := tree.search(req.Path)
		if handler != nil {
			req.Params = params
			return handler(req)
		}
	}

	// 405 check: Does this path exist under ANY other method?
	for _, mTree := range r.trees {
		handler, _ := mTree.search(req.Path)
		if handler != nil {
			return http.NewResponse405()
		}
	}

	// Fallback to 404
	return http.NewResponse404()
}
