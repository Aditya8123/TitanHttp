package middleware

import (
	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// AuthPlaceholder is a simple middleware that checks for the presence of an
// "Authorization" header. If it's missing or incorrect, it returns a 401 Unauthorized.
// This is a placeholder for a future, more robust authentication system.
func AuthPlaceholder(next router.Handler) router.Handler {
	return func(req *http.Request) *http.Response {
		// Get the Authorization header (parser maps headers to lowercase keys)
		authHeader := req.Headers["authorization"]

		// For demonstration, we require a specific hardcoded bearer token
		if authHeader != "Bearer titan-secret-token" {
			return http.NewResponse401()
		}

		// Passed authentication, proceed to the next handler
		return next(req)
	}
}
