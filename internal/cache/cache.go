package cache

import (
	"sync"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

// cacheEntry wraps an HTTP response.
type cacheEntry struct {
	Response  *http.Response
	ExpiresAt time.Time
}

// MemoryCache provides a thread-safe in-memory key-value store for HTTP responses.
type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]cacheEntry
}

// NewMemoryCache initializes and returns a new MemoryCache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]cacheEntry),
	}
}

// StartSweeper begins a background goroutine to periodically clean up expired cache entries.
func (c *MemoryCache) StartSweeper(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			c.mu.Lock()
			now := time.Now()
			for key, entry := range c.store {
				if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
					delete(c.store, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}

// Get retrieves a cached response by key.
func (c *MemoryCache) Get(key string) (*http.Response, bool) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	// Check expiration
	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		// Expired. Lazily delete it.
		c.Delete(key)
		return nil, false
	}

	return entry.Response, true
}

// Set stores a response in the cache with the given key and TTL.
// A TTL of 0 means the entry never expires.
func (c *MemoryCache) Set(key string, resp *http.Response, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if ttl != 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.store[key] = cacheEntry{
		Response:  resp,
		ExpiresAt: expiresAt,
	}
}

// Delete removes a response from the cache by key.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}
