package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestNewRouter(t *testing.T) {
	r, err := NewRouter()
	if err != nil {
		t.Fatalf("NewRouter() returned an error: %v", err)
	}

	if r == nil {
		t.Fatal("NewRouter() returned a nil router")
	}

	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"GET", "/ping"},
		{"GET", "/healthz"},
		{"GET", "/info"},
		{"GET", "/version"},
		{"POST", "/echo"},
	}

	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		var match mux.RouteMatch
		if !r.mux.Match(req, &match) {
			t.Errorf("route not registered: %s %s", route.method, route.path)
		}
	}
}

