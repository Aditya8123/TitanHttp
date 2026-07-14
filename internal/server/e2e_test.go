package server

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	nethttp "net/http"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/cache"
	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/middleware"
	"github.com/Aditya8123/TitanHttp/internal/middleware/rate"
	"github.com/Aditya8123/TitanHttp/internal/proxy"
)

func TestE2E_BasicRequestResponse(t *testing.T) {
	// Spin up a server on dynamic port
	srv := NewServer("127.0.0.1:0")

	srv.Router().Get("/hello", func(req *http.Request) *http.Response {
		resp := http.NewResponse200()
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("Hello, E2E!")
		return resp
	})

	srv.Router().Post("/echo", func(req *http.Request) *http.Response {
		resp := http.NewResponse200()
		resp.Headers.Set("Content-Type", "application/octet-stream")
		resp.Body = req.RawBody
		return resp
	})

	go func() {
		_ = srv.Start()
	}()
	time.Sleep(50 * time.Millisecond)
	addr := srv.Addr()
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	// 1. Test GET
	ctx := context.Background()
	req, err := nethttp.NewRequestWithContext(ctx, "GET", "http://"+addr+"/hello", nethttp.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello, E2E!" {
		t.Errorf("Expected 'Hello, E2E!', got %q", string(body))
	}

	// 2. Test POST
	postReq, err := nethttp.NewRequestWithContext(ctx, "POST", "http://"+addr+"/echo", strings.NewReader("echo payload"))
	if err != nil {
		t.Fatalf("Failed to create post request: %v", err)
	}
	postReq.Header.Set("Content-Type", "text/plain")
	postResp, err := nethttp.DefaultClient.Do(postReq)
	if err != nil {
		t.Fatalf("Failed to POST: %v", err)
	}
	defer postResp.Body.Close()

	if postResp.StatusCode != nethttp.StatusOK {
		t.Errorf("Expected 200, got %d", postResp.StatusCode)
	}

	postBody, _ := io.ReadAll(postResp.Body)
	if string(postBody) != "echo payload" {
		t.Errorf("Expected 'echo payload', got %q", string(postBody))
	}
}

func TestE2E_Middlewares(t *testing.T) {
	srv := NewServer("127.0.0.1:0")

	// Compose with recovery, auth, range, and gzip
	srv.Router().Get("/data", middleware.Recovery(middleware.AuthPlaceholder(middleware.Range(middleware.Gzip(
		func(req *http.Request) *http.Response {
			resp := http.NewResponse200()
			resp.Headers.Set("Content-Type", "text/plain")
			resp.Body = []byte("0123456789 - payload that is sufficiently long for compression check")
			return resp
		},
	)))))

	go func() {
		_ = srv.Start()
	}()
	time.Sleep(50 * time.Millisecond)
	addr := srv.Addr()
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	ctx := context.Background()

	// 1. Unauthorized request
	req, err := nethttp.NewRequestWithContext(ctx, "GET", "http://"+addr+"/data", nethttp.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	client := &nethttp.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != nethttp.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", resp.StatusCode)
	}

	// 2. Successful Authorized request with gzip negotiation
	req2, err := nethttp.NewRequestWithContext(ctx, "GET", "http://"+addr+"/data", nethttp.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req2.Header.Set("Authorization", "Bearer titan-secret-token")
	req2.Header.Set("Accept-Encoding", "gzip")
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != nethttp.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp2.StatusCode)
	}

	if resp2.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("Expected Content-Encoding: gzip, got %q", resp2.Header.Get("Content-Encoding"))
	}

	gr, err := gzip.NewReader(resp2.Body)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	body, _ := io.ReadAll(gr)
	expectedPrefix := "0123456789"
	if !strings.HasPrefix(string(body), expectedPrefix) {
		t.Errorf("Expected body to start with %q, got %q", expectedPrefix, string(body))
	}
}

func TestE2E_ProxyAndLoadBalancer(t *testing.T) {
	// 1. Start target backend server
	backendSrv := NewServer("127.0.0.1:0")
	backendSrv.Router().Get("/info", func(req *http.Request) *http.Response {
		resp := http.NewResponse200()
		resp.Headers.Set("Content-Type", "application/json")
		// Echo back client information injected by proxy
		resp.Body = []byte(fmt.Sprintf(`{"forwarded_for": %q, "forwarded_host": %q}`,
			req.Headers.Get("x-forwarded-for"),
			req.Headers.Get("x-forwarded-host")))
		return resp
	})

	go func() {
		_ = backendSrv.Start()
	}()
	time.Sleep(50 * time.Millisecond)
	backendAddr := backendSrv.Addr()
	defer func() {
		_ = backendSrv.Shutdown(context.Background())
	}()

	// 2. Start proxy load balancer server
	proxySrv := NewServer("127.0.0.1:0")
	lb := proxy.NewLoadBalancer([]string{backendAddr})

	proxySrv.Router().Get("/info", lb)

	go func() {
		_ = proxySrv.Start()
	}()
	time.Sleep(50 * time.Millisecond)
	proxyAddr := proxySrv.Addr()
	defer func() {
		_ = proxySrv.Shutdown(context.Background())
	}()

	// 3. Request through proxy
	ctx := context.Background()
	req, err := nethttp.NewRequestWithContext(ctx, "GET", "http://"+proxyAddr+"/info", nethttp.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to GET through proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "forwarded_for") || !strings.Contains(bodyStr, "127.0.0.1") {
		t.Errorf("Expected injected x-forwarded-for headers, got: %s", bodyStr)
	}
}

func TestE2E_RateLimitingAndCaching(t *testing.T) {
	srv := NewServer("127.0.0.1:0")

	// Limit to 2 requests per second with a bucket size of 2
	limiter := rate.NewTokenBucket(2, 2.0)
	// Cache layer
	memCache := cache.NewMemoryCache()
	cacheLayer := middleware.CacheMiddleware(memCache)

	var requestCount int
	srv.Router().Get("/data", middleware.RateLimitMiddleware(limiter)(cacheLayer(
		func(req *http.Request) *http.Response {
			requestCount++
			resp := http.NewResponse200()
			resp.Headers.Set("Content-Type", "text/plain")
			resp.Headers.Set("Cache-Control", "max-age=5")
			resp.Body = []byte("Counter: " + strconv.Itoa(requestCount))
			return resp
		},
	)))

	go func() {
		_ = srv.Start()
	}()
	time.Sleep(50 * time.Millisecond)
	addr := srv.Addr()
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	client := &nethttp.Client{}
	ctx := context.Background()

	// Request 1: hits the backend, cached
	req1, err := nethttp.NewRequestWithContext(ctx, "GET", "http://"+addr+"/data", nethttp.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request 1: %v", err)
	}
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("Request 1 failed: %v", err)
	}
	body1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()

	if string(body1) != "Counter: 1" {
		t.Errorf("Expected Counter: 1, got %q", string(body1))
	}

	// Request 2: hits the cache, backend counter remains 1
	req2, err := nethttp.NewRequestWithContext(ctx, "GET", "http://"+addr+"/data", nethttp.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request 2: %v", err)
	}
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Request 2 failed: %v", err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()

	if string(body2) != "Counter: 1" {
		t.Errorf("Expected Counter: 1 (cached), got %q", string(body2))
	}

	// Request 3 & 4 (fast): exceeds rate limit -> 429
	for i := 0; i < 3; i++ {
		req, err := nethttp.NewRequestWithContext(ctx, "GET", "http://"+addr+"/data", nethttp.NoBody)
		if err == nil {
			resp, err := client.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == nethttp.StatusTooManyRequests {
					// Hit rate limit successfully
					return
				}
			}
		}
	}
	t.Error("Expected to trigger 429 StatusTooManyRequests rate limiter protection")
}

func TestE2E_HTTPSAndALPN(t *testing.T) {
	// Generate certificates
	tmpDir := t.TempDir()
	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")
	if err := generateTestCertificate(certFile, keyFile); err != nil {
		t.Fatalf("Failed to generate test certificates: %v", err)
	}

	srv := NewServer("127.0.0.1:0")
	srv.Router().Get("/secure", func(req *http.Request) *http.Response {
		resp := http.NewResponse200()
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("Secure Connection")
		return resp
	})

	go func() {
		_ = srv.StartTLS(certFile, keyFile)
	}()
	time.Sleep(100 * time.Millisecond)
	addr := srv.Addr()
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	tr := &nethttp.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &nethttp.Client{Transport: tr}
	ctx := context.Background()

	req, err := nethttp.NewRequestWithContext(ctx, "GET", "https://"+addr+"/secure", nethttp.NoBody)
	if err != nil {
		t.Fatalf("Failed to create HTTPS request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTPS GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Secure Connection" {
		t.Errorf("Expected 'Secure Connection', got %q", string(body))
	}
}
