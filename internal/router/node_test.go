package router

import (
	"reflect"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func dummyHandler(req *http.Request) *http.Response {
	return nil
}

func TestNode_InsertAndSearch(t *testing.T) {
	root := &node{}

	root.insert("/", dummyHandler)
	root.insert("/users", dummyHandler)
	root.insert("/users/:id", dummyHandler)
	root.insert("/users/:id/posts/:post_id", dummyHandler)
	root.insert("/static/index.html", dummyHandler)
	root.insert("/static/*filepath", dummyHandler)
	root.insert("/*catchall", dummyHandler)

	tests := []struct {
		name       string
		searchPath string
		wantMatch  bool
		wantParams map[string]string
	}{
		{
			name:       "Root match",
			searchPath: "/",
			wantMatch:  true,
			wantParams: map[string]string{},
		},
		{
			name:       "Exact match",
			searchPath: "/users",
			wantMatch:  true,
			wantParams: map[string]string{},
		},
		{
			name:       "Parameter match",
			searchPath: "/users/123",
			wantMatch:  true,
			wantParams: map[string]string{"id": "123"},
		},
		{
			name:       "Multi-parameter match",
			searchPath: "/users/123/posts/456",
			wantMatch:  true,
			wantParams: map[string]string{"id": "123", "post_id": "456"},
		},
		{
			name:       "Exact static match",
			searchPath: "/static/index.html",
			wantMatch:  true,
			wantParams: map[string]string{},
		},
		{
			name:       "Wildcard match with segments",
			searchPath: "/static/css/main.css",
			wantMatch:  true,
			wantParams: map[string]string{"filepath": "css/main.css"},
		},
		{
			name:       "Wildcard match single segment",
			searchPath: "/static/logo.png",
			wantMatch:  true,
			wantParams: map[string]string{"filepath": "logo.png"},
		},
		{
			name:       "Wildcard match exact boundary",
			searchPath: "/static/",
			wantMatch:  true,
			wantParams: map[string]string{"filepath": ""},
		},
		{
			name:       "Global catchall",
			searchPath: "/something/completely/different",
			wantMatch:  true,
			wantParams: map[string]string{"catchall": "something/completely/different"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, params := root.search(tt.searchPath)

			if tt.wantMatch && handler == nil {
				t.Errorf("search(%q) expected match, got nil handler", tt.searchPath)
			}
			if !tt.wantMatch && handler != nil {
				t.Errorf("search(%q) expected NO match, got handler", tt.searchPath)
			}

			if tt.wantMatch && !reflect.DeepEqual(params, tt.wantParams) {
				t.Errorf("search(%q) params = %v, want %v", tt.searchPath, params, tt.wantParams)
			}
		})
	}
}

func TestNode_InvalidWildcardPanic(t *testing.T) {
	root := &node{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on invalid wildcard insertion")
		}
	}()

	// This should panic because there are segments after the wildcard
	root.insert("/static/*filepath/invalid", dummyHandler)
}
