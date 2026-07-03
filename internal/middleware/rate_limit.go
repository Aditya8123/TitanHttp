package middleware

import (
	"strings"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/middleware/rate"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// RateLimitMiddleware creates a middleware that enforces rate limits
// using the provided rate.Limiter implementation.
func RateLimitMiddleware(limiter rate.Limiter) router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(req *http.Request) *http.Response {
			// Extract IP from RemoteAddr (e.g. "192.168.1.1:54321")
			ip := req.RemoteAddr
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}

			// In a real application, we would also check X-Forwarded-For if
			// we are behind a proxy. For simplicity, we just use the remote IP.
			if forwardedFor := req.Headers["x-forwarded-for"]; forwardedFor != "" {
				// Get the first IP in the list
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					ip = strings.TrimSpace(ips[0])
				}
			}

			if !limiter.Allow(ip) {
				return http.NewResponse429()
			}

			return next(req)
		}
	}
}
