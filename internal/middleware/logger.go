package middleware

import (
	"fmt"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// Logger is a middleware that logs incoming HTTP requests and the resulting
// HTTP response status codes along with the request latency.
func Logger(next router.Handler) router.Handler {
	return func(req *http.Request) *http.Response {
		start := time.Now()

		// Pre-processing: log the incoming request
		fmt.Printf("[TitanHTTP] → %s %s %s\n", req.Method, req.Path, req.Version)

		// Pass execution to the next handler
		resp := next(req)

		// Post-processing: log the response details and latency
		latency := time.Since(start)
		fmt.Printf("[TitanHTTP] ← %d %s (%v)\n", resp.StatusCode, http.StatusText(resp.StatusCode), latency)

		return resp
	}
}
