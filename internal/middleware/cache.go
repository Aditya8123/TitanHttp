package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/cache"
	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// CacheMiddleware creates a middleware that intercepts requests and serves
// them from the provided cache if a fresh entry exists. If not, it forwards
// the request to the next handler and caches the resulting response.
func CacheMiddleware(c *cache.MemoryCache) router.Middleware {
	return func(next router.Handler) router.Handler {
		return func(req *http.Request) *http.Response {
			// Only cache GET requests.
			if req.Method != http.MethodGet {
				return next(req)
			}

			reqCC := req.Headers["cache-control"]
			bypassCache := strings.Contains(reqCC, "no-cache") || strings.Contains(reqCC, "no-store")

			if !bypassCache {
				// Try to serve from cache
				cacheKey := req.Path
				if cachedResp, ok := c.Get(cacheKey); ok {
					return cachedResp
				}
			}

			// Cache miss or bypass, execute next handler
			resp := next(req)

			// Only cache 200 OK responses and fully buffered bodies
			if resp.StatusCode == 200 && resp.Stream == nil {
				respCC := resp.Headers["cache-control"]
				if !strings.Contains(respCC, "no-store") && !strings.Contains(respCC, "no-cache") {
					// Parse max-age if present
					ttl := time.Duration(0) // Default: never expires
					if strings.Contains(respCC, "max-age=") {
						parts := strings.Split(respCC, "max-age=")
						if len(parts) > 1 {
							// Extract just the number
							valStr := strings.Split(parts[1], ",")[0]
							valStr = strings.TrimSpace(valStr)
							if secs, err := strconv.Atoi(valStr); err == nil && secs > 0 {
								ttl = time.Duration(secs) * time.Second
							}
						}
					}
					c.Set(req.Path, resp, ttl)
				}
			}

			return resp
		}
	}
}
