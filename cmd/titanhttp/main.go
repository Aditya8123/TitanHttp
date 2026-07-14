package main

import (
	"fmt"
	"log"
	nethttp "net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/middleware"
	"github.com/Aditya8123/TitanHttp/internal/server"
)

func main() {
	fmt.Println("Initializing TitanHTTP Server...")

	// Start pprof debug server in background
	go func() {
		log.Println("Starting pprof debug server on localhost:6060")
		debugSrv := &nethttp.Server{
			Addr:         "localhost:6060",
			Handler:      nil,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		log.Println(debugSrv.ListenAndServe())
	}()

	port := os.Getenv("PORT")
	if port == "" {
		if os.Getenv("TLS_CERT") != "" {
			port = "8443"
		} else {
			port = "8080"
		}
	}

	srv := server.NewServer(":" + port)
	if timeoutEnv := os.Getenv("IDLE_TIMEOUT"); timeoutEnv != "" {
		if d, err := time.ParseDuration(timeoutEnv); err == nil {
			srv.IdleTimeout = d
			fmt.Printf("Configured idle timeout from environment: %v\n", d)
		}
	}

	// Mount global middlewares
	// srv.Router().Use(middleware.Logger) // Disabled for benchmarking to avoid stdout mutex blocking
	srv.Router().Use(middleware.Recovery)
	srv.Router().Use(middleware.Range)
	srv.Router().Use(middleware.Gzip)

	// Register some basic routes to demonstrate the new Router
	srv.Router().Get("/", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("Welcome to TitanHTTP!\n")
		return resp
	})

	srv.Router().Get("/hello", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("Hello from the new Router!\n")
		return resp
	})

	srv.Router().Get("/ping", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("pong")
		return resp
	})

	srv.Router().Post("/json", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "application/json")
		resp.Body = []byte(`{"status":"success"}`)
		return resp
	})

	// Pre-generate a 100KB payload for the heavy route to accurately test network I/O
	// instead of memory allocation bottlenecks.
	heavyPayload := []byte(strings.Repeat("junk_data_", 10000))

	srv.Router().Get("/heavy", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = heavyPayload
		return resp
	})

	srv.Router().Get("/users/:name", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = fmt.Appendf(nil, "Hello, %s!\n", req.Params["name"])
		return resp
	})

	srv.Router().Get("/static/*filepath", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		filepath := "public/" + req.Params["filepath"]
		data, err := os.ReadFile(filepath)
		if err != nil {
			resp.StatusCode = http.StatusNotFound
			resp.Headers.Set("Content-Type", "text/plain")
			resp.Body = []byte("404 File Not Found\n")
			return resp
		}
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = data
		return resp
	})

	srv.Router().Get("/panic", func(req *http.Request) *http.Response {
		panic("This is a simulated panic!")
	})

	srv.Router().Get("/protected", middleware.AuthPlaceholder(func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("Welcome to the secret protected area!\n")
		return resp
	}))

	// --- Benchmark Additions ---

	// Payload matrix endpoint
	srv.Router().Get("/payload/:size", func(req *http.Request) *http.Response {
		resp := http.NewResponse()

		sizeStr := req.Params["size"]
		var size int
		fmt.Sscanf(sizeStr, "%d", &size)

		if size < 0 || size > 15_000_000 {
			resp.StatusCode = http.StatusBadRequest
			resp.Headers.Set("Content-Type", "text/plain")
			resp.Body = []byte("Size must be between 0 and 10MB")
			return resp
		}

		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "application/octet-stream")
		resp.Body = make([]byte, size)
		return resp
	})

	// POST endpoint
	srv.Router().Post("/api/users", func(req *http.Request) *http.Response {
		// Simulating reading body and unmarshaling. We don't have json.Unmarshal right here,
		// but we can parse the body simulating standard work.
		// For benchmark purposes, we just return the {"ok":true} directly.
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "application/json")
		resp.Body = []byte(`{"ok":true}`)
		return resp
	})

	// GC endpoint for memory profiling stage
	srv.Router().Get("/debug/gc", func(req *http.Request) *http.Response {
		runtime.GC()
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("GC complete")
		return resp
	})

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
