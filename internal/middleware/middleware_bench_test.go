package middleware

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/middleware/rate"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

func init() {
	LogOutput = io.Discard
}

func dummyHandler(req *http.Request) *http.Response {
	res := http.NewResponse()
	res.StatusCode = http.StatusOK
	res.Body = []byte("OK")
	return res
}

func setupReq() *http.Request {
	reqStr := "GET / HTTP/1.1\r\nHost: localhost\r\nAuthorization: Bearer secret-token\r\nAccept-Encoding: gzip\r\n\r\n"
	req, _ := http.ParseRequest(bufio.NewReader(bytes.NewReader([]byte(reqStr))))
	req.RemoteAddr = "127.0.0.1:12345"
	return req
}

func BenchmarkMiddleware_None(b *testing.B) {
	req := setupReq()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := dummyHandler(req)
		_ = res.Body
	}
}

func BenchmarkMiddleware_Logger(b *testing.B) {
	req := setupReq()
	chain := router.Chain(Logger)(dummyHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := chain(req)
		_ = res.Body
	}
}

func BenchmarkMiddleware_LoggerRecovery(b *testing.B) {
	req := setupReq()
	chain := router.Chain(Logger, Recovery)(dummyHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := chain(req)
		_ = res.Body
	}
}

func BenchmarkMiddleware_LoggerRecoveryAuth(b *testing.B) {
	req := setupReq()
	chain := router.Chain(Logger, Recovery, AuthPlaceholder)(dummyHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := chain(req)
		_ = res.Body
	}
}

func BenchmarkMiddleware_LoggerRecoveryAuthRateLimit(b *testing.B) {
	req := setupReq()
	tb := rate.NewTokenBucket(10000, 10000)
	chain := router.Chain(Logger, Recovery, AuthPlaceholder, RateLimitMiddleware(tb))(dummyHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := chain(req)
		_ = res.Body
	}
}

func BenchmarkMiddleware_LoggerRecoveryAuthRateLimitGzip(b *testing.B) {
	req := setupReq()
	tb := rate.NewTokenBucket(10000, 10000)
	chain := router.Chain(Logger, Recovery, AuthPlaceholder, RateLimitMiddleware(tb), Gzip)(dummyHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := chain(req)
		_ = res.Body
	}
}
