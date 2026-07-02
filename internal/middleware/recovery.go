package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// Recovery is a middleware that recovers from any panics in subsequent handlers
// and returns a 500 Internal Server Error response, preventing the server from crashing.
func Recovery(next router.Handler) router.Handler {
	return func(req *http.Request) (resp *http.Response) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				fmt.Printf("[TitanHTTP] 💥 PANIC RECOVERED: %v\n", err)
				debug.PrintStack()

				// Return a 500 Internal Server Error
				resp = http.NewResponse500()
			}
		}()

		// Call the next handler
		resp = next(req)
		return resp
	}
}
