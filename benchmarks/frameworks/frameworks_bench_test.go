package frameworks_bench

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/go-chi/chi/v5"

	thttp "github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func BenchmarkNetHTTP(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

func BenchmarkGin(b *testing.B) {
	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

func BenchmarkChi(b *testing.B) {
	router := chi.NewRouter()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
		w.Body.Reset()
	}
}

func BenchmarkFiber(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.Test(req, -1)
	}
}

func BenchmarkTitanHTTP(b *testing.B) {
	r := router.NewRouter()
	r.Get("/ping", func(req *thttp.Request) *thttp.Response {
		resp := thttp.NewResponse()
		resp.StatusCode = thttp.StatusOK
		resp.Body = []byte("pong")
		return resp
	})

	req := thttp.NewRequest()
	req.Method = thttp.MethodGet
	req.Path = "/ping"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resp := r.ServeHTTP(req)
		thttp.ReleaseResponse(resp)
	}
}
