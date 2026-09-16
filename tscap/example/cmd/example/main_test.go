package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	goahttp "goa.design/goa/v3/http"
	"goa.design/plugins/v3/tscap/example"
	tscapsvr "goa.design/plugins/v3/tscap/example/gen/http/tscap/server"
	tscapsvc "goa.design/plugins/v3/tscap/example/gen/tscap"
)

func setupServer() http.Handler {
	logger := log.New(os.Stderr, "[test] ", log.Ltime)
	svc := example.NewTscap(logger)
	endpoints := tscapsvc.NewEndpoints(svc)
	mux := goahttp.NewMuxer()
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	svr := tscapsvr.New(endpoints, mux, dec, enc, nil, nil)
	tscapsvr.Mount(mux, svr)
	return mux
}

func makeCapHeader(caps map[string][]map[string][]string) string {
	b, _ := json.Marshal(caps)
	return string(b)
}

func TestHealthAnonymous(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestListWithoutCaps(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestListWithReadCap(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("GET", "/items", nil)
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/tscap": {
			{"action": {"read"}, "resources": {"*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestListWithWildcardCap(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("GET", "/items", nil)
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/tscap": {
			{"action": {"*"}, "resources": {"*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestListWithWrongCap(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("GET", "/items", nil)
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/tscap": {
			{"action": {"write"}, "resources": {"*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestListWithWrongCapability(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("GET", "/items", nil)
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/other": {
			{"action": {"read"}, "resources": {"*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestCreateWithWriteCap(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("POST", "/items", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/tscap": {
			{"action": {"write"}, "resources": {"items/*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		body, _ := io.ReadAll(w.Body)
		t.Errorf("expected 201, got %d: %s", w.Code, string(body))
	}
}

func TestCreateWithReadCap(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("POST", "/items", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/tscap": {
			{"action": {"read"}, "resources": {"*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestAdminWithAdminCap(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("DELETE", "/items/123", nil)
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/tscap": {
			{"action": {"admin"}, "resources": {"*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		body, _ := io.ReadAll(w.Body)
		t.Errorf("expected 204, got %d: %s", w.Code, string(body))
	}
}

func TestAdminWithReadCap(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest("DELETE", "/items/123", nil)
	req.Header.Set("Tailscale-App-Capabilities", makeCapHeader(map[string][]map[string][]string{
		"example.com/cap/tscap": {
			{"action": {"read"}, "resources": {"*"}},
		},
	}))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
