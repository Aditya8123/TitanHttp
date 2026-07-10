package middleware

import (
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// LogOutput allows redirecting the logger's output (defaults to os.Stdout).
var LogOutput io.Writer = os.Stdout

var logBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 128)
		return &b
	},
}

// Logger is a middleware that logs incoming HTTP requests and the resulting
// HTTP response status codes along with the request latency.
func Logger(next router.Handler) router.Handler {
	return func(req *http.Request) *http.Response {
		start := time.Now()

		// Pre-processing: log the incoming request
		ptr := logBufPool.Get().(*[]byte)
		buf := (*ptr)[:0]

		buf = append(buf, "[TitanHTTP] → "...)
		buf = append(buf, req.Method...)
		buf = append(buf, ' ')
		buf = append(buf, req.Path...)
		buf = append(buf, ' ')
		buf = append(buf, req.Version...)
		buf = append(buf, '\n')
		LogOutput.Write(buf)

		// Pass execution to the next handler
		resp := next(req)

		// Post-processing: log the response details and latency
		buf = buf[:0]
		buf = append(buf, "[TitanHTTP] ← "...)
		buf = strconv.AppendInt(buf, int64(resp.StatusCode), 10)
		buf = append(buf, ' ')
		buf = append(buf, http.StatusText(resp.StatusCode)...)
		buf = append(buf, " ("...)
		buf = append(buf, time.Since(start).String()...)
		buf = append(buf, ")\n"...)
		LogOutput.Write(buf)

		*ptr = buf
		logBufPool.Put(ptr)

		return resp
	}
}
