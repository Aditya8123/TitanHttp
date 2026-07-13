package middleware

import (
	"strconv"
	"strings"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// Range handles basic Range requests, returning 206 Partial Content.
// For a fully compliant server, this would slice the response body stream.
// Since TitanHTTP buffers body responses currently, we slice the byte array.
func Range(next router.Handler) router.Handler {
	return func(req *http.Request) *http.Response {
		resp := next(req)

		rangeHeader := req.Headers.Get("range")
		if rangeHeader == "" || resp.StatusCode != http.StatusOK || len(resp.Body) == 0 {
			return resp
		}

		if strings.HasPrefix(rangeHeader, "bytes=") {
			// Basic implementation for bytes=X-Y
			parts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
			if len(parts) == 2 {
				start, err1 := strconv.Atoi(parts[0])
				end, err2 := strconv.Atoi(parts[1])
				
				if err1 == nil && err2 == nil && start >= 0 && end >= start && end < len(resp.Body) {
					resp.StatusCode = http.StatusPartialContent
					
					origLen := len(resp.Body)
					resp.Body = resp.Body[start : end+1]
					
					resp.Headers.Set("Content-Range", "bytes "+strconv.Itoa(start)+"-"+strconv.Itoa(end)+"/"+strconv.Itoa(origLen))
					resp.Headers.Set("Content-Length", strconv.Itoa(len(resp.Body)))
				}
			}
		}

		return resp
	}
}
