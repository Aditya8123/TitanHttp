package http

import (
	"strings"
)

type headerEntry struct {
	key   string
	value string
}

// Header represents HTTP headers backed by a slice for zero-allocation parsing.
// It uses a lazy index for fast lookups when the number of headers exceeds a threshold.
type Header struct {
	entries []headerEntry
	index   map[string][]int
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
func (h *Header) Add(key, value string) {
	k := strings.ToLower(key)
	h.entries = append(h.entries, headerEntry{key: k, value: value})
	if h.index != nil {
		h.index[k] = append(h.index[k], len(h.entries)-1)
	} else if len(h.entries) > 16 {
		h.buildIndex()
	}
}

// Set sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key.
func (h *Header) Set(key, value string) {
	k := strings.ToLower(key)
	
	// If index is built, we can just replace the first one and remove others,
	// or rebuild. For simplicity, we just delete and add.
	h.Del(k)
	h.Add(k, value)
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
func (h *Header) Get(key string) string {
	k := strings.ToLower(key)
	if h.index != nil {
		if indices, ok := h.index[k]; ok && len(indices) > 0 {
			return h.entries[indices[0]].value
		}
		return ""
	}

	for _, entry := range h.entries {
		if entry.key == k {
			return entry.value
		}
	}
	return ""
}

// GetAll returns all values associated with the given key.
func (h *Header) GetAll(key string) []string {
	k := strings.ToLower(key)
	var values []string
	if h.index != nil {
		if indices, ok := h.index[k]; ok {
			for _, idx := range indices {
				values = append(values, h.entries[idx].value)
			}
		}
		return values
	}

	for _, entry := range h.entries {
		if entry.key == k {
			values = append(values, entry.value)
		}
	}
	return values
}

// Del deletes the values associated with key.
func (h *Header) Del(key string) {
	k := strings.ToLower(key)
	
	// Fast path: rebuild entries without the key
	var newEntries []headerEntry
	for _, entry := range h.entries {
		if entry.key != k {
			newEntries = append(newEntries, entry)
		}
	}
	h.entries = newEntries
	
	if h.index != nil {
		h.buildIndex()
	}
}

func (h *Header) buildIndex() {
	h.index = make(map[string][]int)
	for i, entry := range h.entries {
		h.index[entry.key] = append(h.index[entry.key], i)
	}
}

// Reset clears the headers for pool reuse.
func (h *Header) Reset() {
	h.entries = h.entries[:0]
	h.index = nil
}

// Entries allows iterating over all headers.
func (h *Header) Entries() []headerEntry {
	return h.entries
}

// Key returns the header key.
func (e headerEntry) Key() string {
	return e.key
}

// Value returns the header value.
func (e headerEntry) Value() string {
	return e.value
}
