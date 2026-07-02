package main

import (
	"fmt"
	"os"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
	"github.com/Aditya8123/TitanHttp/internal/server"
)

func main() {
	fmt.Println("Initializing TitanHTTP Server...")

	srv := server.NewServer(":8080")

	// Add a dummy global middleware to prove the pipeline works
	srv.Router().Use(func(next router.Handler) router.Handler {
		return func(req *http.Request) *http.Response {
			fmt.Printf("[Middleware] Intercepted %s %s\n", req.Method, req.Path)
			
			// Call the next handler in the chain
			resp := next(req)
			
			// Post-process the response
			resp.Headers["X-Titan-Middleware"] = "Pipeline Active"
			return resp
		}
	})

	// Register some basic routes to demonstrate the new Router
	srv.Router().Get("/", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/plain"
		resp.Body = []byte("Welcome to TitanHTTP!\n")
		return resp
	})

	srv.Router().Get("/hello", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/plain"
		resp.Body = []byte("Hello from the new Router!\n")
		return resp
	})

	srv.Router().Get("/users/:name", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/plain"
		resp.Body = []byte(fmt.Sprintf("Hello, %s!\n", req.Params["name"]))
		return resp
	})

	srv.Router().Get("/static/*filepath", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/plain"
		resp.Body = []byte(fmt.Sprintf("Serving static file: %s\n", req.Params["filepath"]))
		return resp
	})

	if err := srv.Start(); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}
