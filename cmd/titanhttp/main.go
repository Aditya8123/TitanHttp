package main

import (
	"fmt"
	"log"
	nethttp "net/http"
	_ "net/http/pprof"
	"os"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/middleware"
	"github.com/Aditya8123/TitanHttp/internal/server"
)

func main() {
	fmt.Println("Initializing TitanHTTP Server...")

	// Start pprof debug server in background
	go func() {
		log.Println("Starting pprof debug server on localhost:6060")
		log.Println(nethttp.ListenAndServe("localhost:6060", nil))
	}()

	srv := server.NewServer(":8080")

	// Mount global middlewares
	srv.Router().Use(middleware.Logger)
	srv.Router().Use(middleware.Recovery)
	srv.Router().Use(middleware.Gzip)

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

	srv.Router().Get("/ping", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/plain"
		resp.Body = []byte("pong")
		return resp
	})

	srv.Router().Get("/heavy", func(req *http.Request) *http.Response {
		// Simulate a memory-heavy allocation to test GC pressure
		var data []byte
		for i := 0; i < 10000; i++ {
			data = append(data, []byte("junk_data_")...)
		}
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/plain"
		resp.Body = []byte(fmt.Sprintf("Allocated %d bytes", len(data)))
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

	srv.Router().Get("/panic", func(req *http.Request) *http.Response {
		panic("This is a simulated panic!")
	})

	srv.Router().Get("/protected", middleware.AuthPlaceholder(func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/plain"
		resp.Body = []byte("Welcome to the secret protected area!\n")
		return resp
	}))

	certFile := os.Getenv("TLS_CERT")
	keyFile := os.Getenv("TLS_KEY")

	if certFile != "" && keyFile != "" {
		fmt.Printf("Starting in HTTPS mode using cert: %s\n", certFile)
		if err := srv.StartTLS(certFile, keyFile); err != nil {
			fmt.Printf("Fatal TLS error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := srv.Start(); err != nil {
			fmt.Printf("Fatal error: %v\n", err)
			os.Exit(1)
		}
	}
}
