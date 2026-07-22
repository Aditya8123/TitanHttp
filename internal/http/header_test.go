package http

import (
	"testing"
)

func TestHeaderBasic(t *testing.T) {
	var h Header

	h.Add("Content-Type", "application/json")
	h.Add("Host", "localhost")

	if val := h.Get("Content-Type"); val != "application/json" {
		t.Errorf("Expected application/json, got %q", val)
	}
	if val := h.Get("content-type"); val != "application/json" {
		t.Errorf("Expected case-insensitive get to work, got %q", val)
	}
	if val := h.Get("Host"); val != "localhost" {
		t.Errorf("Expected localhost, got %q", val)
	}

	h.Set("Content-Type", "text/plain")
	if val := h.Get("Content-Type"); val != "text/plain" {
		t.Errorf("Expected text/plain after Set, got %q", val)
	}

	h.Del("Host")
	if val := h.Get("Host"); val != "" {
		t.Errorf("Expected Host to be deleted, got %q", val)
	}
}

func TestHeaderMultipleValues(t *testing.T) {
	var h Header

	h.Add("Accept", "text/html")
	h.Add("Accept", "application/xhtml+xml")
	h.Add("Accept", "application/xml")

	vals := h.GetAll("Accept")
	if len(vals) != 3 {
		t.Fatalf("Expected 3 values, got %d", len(vals))
	}
	if vals[0] != "text/html" || vals[1] != "application/xhtml+xml" || vals[2] != "application/xml" {
		t.Errorf("Unexpected values: %v", vals)
	}

	// GetAll on non-existent key
	if len(h.GetAll("Non-Existent")) != 0 {
		t.Errorf("Expected empty slice for non-existent key")
	}
}

func TestHeaderLazyIndexing(t *testing.T) {
	var h Header

	// Add more than 16 headers to trigger index build
	for i := 0; i < 20; i++ {
		h.Add(string(rune('A'+i)), "val")
	}

	if h.index == nil {
		t.Errorf("Expected index to be built")
	}

	// Verify lookups using index
	if val := h.Get("A"); val != "val" {
		t.Errorf("Expected A to be val, got %q", val)
	}

	// Reset should clear index and entries
	h.Reset()
	if len(h.entries) != 0 || h.index != nil {
		t.Errorf("Expected Reset to clear all fields")
	}
}

func TestHeaderEntries(t *testing.T) {
	var h Header
	h.Add("Key1", "Value1")
	h.Add("Key2", "Value2")

	entries := h.Entries()
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	if entries[0].Key() != "key1" || entries[0].Value() != "Value1" {
		t.Errorf("Unexpected entry at 0: %s = %s", entries[0].Key(), entries[0].Value())
	}
	if entries[1].Key() != "key2" || entries[1].Value() != "Value2" {
		t.Errorf("Unexpected entry at 1: %s = %s", entries[1].Key(), entries[1].Value())
	}
}
