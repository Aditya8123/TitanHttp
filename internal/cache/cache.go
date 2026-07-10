package cache

import (
	"hash/fnv"
	"sync"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

const numShards = 32

// cacheEntry wraps an HTTP response.
type cacheEntry struct {
	Response  *http.Response
	ExpiresAt time.Time
}

// cacheShard is a single bucket in the memory cache.
type cacheShard struct {
	mu    sync.RWMutex
	store map[string]cacheEntry
}

// MemoryCache provides a thread-safe in-memory key-value store for HTTP responses.
// It uses sharding to reduce lock contention under high concurrency.
type MemoryCache struct {
	shards [numShards]*cacheShard
}

// NewMemoryCache initializes and returns a new MemoryCache.
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{}
	for i := 0; i < numShards; i++ {
		c.shards[i] = &cacheShard{
			store: make(map[string]cacheEntry),
		}
	}
	return c
}

func (c *MemoryCache) getShard(key string) *cacheShard {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	return c.shards[hasher.Sum32()%numShards]
}

// StartSweeper begins a background goroutine to periodically clean up expired cache entries.
func (c *MemoryCache) StartSweeper(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			for i := 0; i < numShards; i++ {
				shard := c.shards[i]
				shard.mu.Lock()
				for key, entry := range shard.store {
					if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
						delete(shard.store, key)
					}
				}
				shard.mu.Unlock()
			}
		}
	}()
}

// Get retrieves a cached response by key.
func (c *MemoryCache) Get(key string) (*http.Response, bool) {
	shard := c.getShard(key)
	shard.mu.RLock()
	entry, ok := shard.store[key]
	shard.mu.RUnlock()

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
	var expiresAt time.Time
	if ttl != 0 {
		expiresAt = time.Now().Add(ttl)
	}

	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.store[key] = cacheEntry{
		Response:  resp,
		ExpiresAt: expiresAt,
	}
}

// Delete removes a response from the cache by key.
func (c *MemoryCache) Delete(key string) {
	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.store, key)
}
