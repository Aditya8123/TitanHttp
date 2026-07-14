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

			reqCC := req.Headers.Get("cache-control")
			bypassCache := strings.Contains(reqCC, "no-cache") || strings.Contains(reqCC, "no-store")

			if !bypassCache {
				// Try to serve from cache
				cacheKey := req.Path
				if cachedResp, ok := c.Get(cacheKey); ok {
					// Acquire a response from the pool and clone the cached fields into it.
					// The server will release this acquired response when done.
					resp := http.AcquireResponse()
					resp.StatusCode = cachedResp.StatusCode
					resp.StatusText = cachedResp.StatusText
					resp.Version = cachedResp.Version

					for _, entry := range cachedResp.Headers.Entries() {
						resp.Headers.Add(entry.Key(), entry.Value())
					}

					if cachedResp.Body != nil {
						resp.Body = make([]byte, len(cachedResp.Body))
						copy(resp.Body, cachedResp.Body)
					}
					return resp
				}
			}

			// Cache miss or bypass, execute next handler
			resp := next(req)

			// Only cache 200 OK responses and fully buffered bodies
			if resp.StatusCode == 200 && resp.Stream == nil {
				respCC := resp.Headers.Get("cache-control")
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
					// Create a decoupled clone of the response for cache storage.
					// This prevents it from being zeroed out when the server calls ReleaseResponse.
					cachedResp := &http.Response{
						StatusCode: resp.StatusCode,
						StatusText: resp.StatusText,
						Version:    resp.Version,
					}
					for _, entry := range resp.Headers.Entries() {
						cachedResp.Headers.Add(entry.Key(), entry.Value())
					}
					if resp.Body != nil {
						cachedResp.Body = make([]byte, len(resp.Body))
						copy(cachedResp.Body, resp.Body)
					}
					c.Set(req.Path, cachedResp, ttl)
				}
			}

			return resp
		}
	}
}
