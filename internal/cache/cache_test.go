package cache

import (
	"testing"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestMemoryCache_GetSet(t *testing.T) {
	c := NewMemoryCache()

	// Miss
	if _, ok := c.Get("/foo"); ok {
		t.Error("Expected cache miss for /foo")
	}

	// Set
	resp := http.NewResponse200()
	resp.Body = []byte("Hello Cache")
	c.Set("/foo", resp, 0)

	// Hit
	hit, ok := c.Get("/foo")
	if !ok {
		t.Fatal("Expected cache hit for /foo")
	}
	if string(hit.Body) != "Hello Cache" {
		t.Errorf("Expected body 'Hello Cache', got '%s'", string(hit.Body))
	}

	// Delete
	c.Delete("/foo")
	if _, ok := c.Get("/foo"); ok {
		t.Error("Expected cache miss after delete")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := NewMemoryCache()

	resp := http.NewResponse200()
	// TTL of -1 millisecond means it expires immediately
	c.Set("/expired", resp, -1 * time.Millisecond)

	if _, ok := c.Get("/expired"); ok {
		t.Error("Expected cache miss for immediately expired item")
	}
}
