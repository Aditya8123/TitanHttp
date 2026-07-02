package main

import (
	"fmt"
	"os"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/server"
)

func main() {
	fmt.Println("Initializing TitanHTTP Server...")

	srv := server.NewServer(":8080")

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

	if err := srv.Start(); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}
