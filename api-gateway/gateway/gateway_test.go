package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-gateway/config"
)

func newTestGateway(t *testing.T) http.Handler {
	t.Helper()

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend got " + r.URL.Path))
	}))
	t.Cleanup(backend.Close)

	h, err := New([]config.Route{{Prefix: "/users", Target: backend.URL}}, "secret")
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestHealthIsOpen(t *testing.T) {
	h := newTestGateway(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
}

func TestRejectsMissingKey(t *testing.T) {
	h := newTestGateway(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/users/42", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want 401", rec.Code)
	}
}

func TestProxiesWithKeyAndStripsPrefix(t *testing.T) {
	h := newTestGateway(t)
	req := httptest.NewRequest("GET", "/users/42", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
	if got := rec.Body.String(); got != "backend got /42" {
		t.Fatalf("got %q, want %q", got, "backend got /42")
	}
}
