package router

import (
	"sync"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

// Handler defines the function signature for processing HTTP requests.
type Handler func(req *http.Request) *http.Response

// Router manages the registration and dispatching of HTTP routes using Radix trees.
type Router struct {
	// mu protects the trees and handler during concurrent read/write operations
	mu          sync.RWMutex
	trees       map[http.Method]*node
	middlewares []Middleware
	handler     Handler
}

// NewRouter initializes and returns a new Router instance.
func NewRouter() *Router {
	r := &Router{
		trees: make(map[http.Method]*node),
	}
	r.handler = r.serveHTTP
	return r
}

// Use adds global middlewares to the router.
// Middlewares are executed in the order they are added.
func (r *Router) Use(middlewares ...Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, middlewares...)
	r.handler = Chain(r.middlewares...)(r.serveHTTP)
}

// AddRoute registers a new handler for the given HTTP method and pattern.
func (r *Router) AddRoute(method http.Method, pattern string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
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

// ServeHTTP processes the request by executing the pre-compiled middleware chain.
func (r *Router) ServeHTTP(req *http.Request) *http.Response {
	r.mu.RLock()
	h := r.handler
	r.mu.RUnlock()
	return h(req)
}

// serveHTTP is the core routing logic that matches paths and extracts parameters.
func (r *Router) serveHTTP(req *http.Request) *http.Response {
	if req.Params == nil {
		req.Params = make(map[string]string)
	}

	r.mu.RLock()
	tree, methodExists := r.trees[req.Method]
	var handler Handler
	if methodExists {
		handler = tree.search(req.Path, req.Params)
	}
	r.mu.RUnlock()

	if handler != nil {
		return handler(req)
	}

	// 405 check: Does this path exist under ANY other method?
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, mTree := range r.trees {
		handler := mTree.search(req.Path, req.Params)
		if handler != nil {
			return http.NewResponse405()
		}
	}

	// Fallback to 404
	return http.NewResponse404()
}

