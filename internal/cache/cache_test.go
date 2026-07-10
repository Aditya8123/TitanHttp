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

// --- Benchmarks ---

func BenchmarkCache_Hit(b *testing.B) {
	c := NewMemoryCache()
	resp := http.NewResponse200()
	c.Set("/popular", resp, 0)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = c.Get("/popular")
	}
}

func BenchmarkCache_Miss(b *testing.B) {
	c := NewMemoryCache()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = c.Get("/missing")
	}
}

func BenchmarkCache_Parallel(b *testing.B) {
	c := NewMemoryCache()
	resp := http.NewResponse200()
	// Create 100 popular routes to test shard distribution
	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		keys[i] = "/popular/" + string(rune('a'+i%26))
		c.Set(keys[i], resp, 0)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = c.Get(keys[i%100])
			i++
		}
	})
}
