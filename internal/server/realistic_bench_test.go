package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/router"
	"github.com/Aditya8123/TitanHttp/internal/middleware"
	titanhttp "github.com/Aditya8123/TitanHttp/internal/http"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func setupRealisticRouter() *router.Router {
	r := router.NewRouter()

	// Apply multiple middlewares
	r.Use(middleware.Recovery)
	// Skipping logger for benchmark to avoid noisy stdout, though real-world would have it
	// We will benchmark middleware overhead separately.

	r.Get("/products", func(req *titanhttp.Request) *titanhttp.Response {
		products := []Product{{ID: 1, Name: "Tablet", Price: 299.99}, {ID: 2, Name: "Phone", Price: 599.99}}
		b, _ := json.Marshal(products)
		res := titanhttp.NewResponse()
		res.StatusCode = titanhttp.StatusOK
		res.Headers.Set("Content-Type", "application/json")
		res.Body = b
		return res
	})

	r.Get("/products/:id", func(req *titanhttp.Request) *titanhttp.Response {
		p := Product{ID: 15, Name: "Tablet", Price: 299.99}
		b, _ := json.Marshal(p)
		res := titanhttp.NewResponse()
		res.StatusCode = titanhttp.StatusOK
		res.Headers.Set("Content-Type", "application/json")
		res.Body = b
		return res
	})

	r.Post("/products", func(req *titanhttp.Request) *titanhttp.Response {
		// Simulate reading and parsing JSON
		var p Product
		bodyBytes := req.Body
		json.Unmarshal(bodyBytes, &p)

		res := titanhttp.NewResponse()
		res.StatusCode = titanhttp.StatusCreated
		res.Headers.Set("Content-Type", "application/json")
		// Echo back
		res.Body = bodyBytes
		return res
	})

	return r
}

func BenchmarkRealisticAPI_GetList(b *testing.B) {
	r := setupRealisticRouter()
	reqBytes := []byte("GET /products HTTP/1.1\r\nHost: localhost\r\n\r\n")

	bytesReader := bytes.NewReader(reqBytes)
	bufReader := bufio.NewReader(bytesReader)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytesReader.Reset(reqBytes)
		bufReader.Reset(bytesReader)

		req, _ := titanhttp.ParseRequest(bufReader)
		res := r.ServeHTTP(req)
		
		// Consume body
		_ = res.Body
		
		titanhttp.ReleaseRequest(req)
		titanhttp.ReleaseResponse(res)
	}
}

func BenchmarkRealisticAPI_GetParam(b *testing.B) {
	r := setupRealisticRouter()
	reqBytes := []byte("GET /products/15 HTTP/1.1\r\nHost: localhost\r\n\r\n")

	bytesReader := bytes.NewReader(reqBytes)
	bufReader := bufio.NewReader(bytesReader)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytesReader.Reset(reqBytes)
		bufReader.Reset(bytesReader)

		req, err := titanhttp.ParseRequest(bufReader)
		if err != nil {
			b.Fatalf("ParseRequest error: %v", err)
		}
		res := r.ServeHTTP(req)
		
		_ = res.Body
		
		titanhttp.ReleaseRequest(req)
		titanhttp.ReleaseResponse(res)
	}
}

func BenchmarkRealisticAPI_PostJSON(b *testing.B) {
	r := setupRealisticRouter()
	body := `{"id":15,"name":"Tablet","price":299.99}`
	reqBytes := []byte("POST /products HTTP/1.1\r\nHost: localhost\r\nContent-Length: 39\r\n\r\n" + body)
	
	bytesReader := bytes.NewReader(reqBytes)
	bufReader := bufio.NewReader(bytesReader)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytesReader.Reset(reqBytes)
		bufReader.Reset(bytesReader)

		req, _ := titanhttp.ParseRequest(bufReader)
		res := r.ServeHTTP(req)
		
		_ = res.Body
		
		titanhttp.ReleaseRequest(req)
		titanhttp.ReleaseResponse(res)
	}
}
