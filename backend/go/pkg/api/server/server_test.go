//go:build apitests

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	r := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestOpenAPI(t *testing.T) {
	r := New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/openapi.json", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	ct := w.Result().Header.Get("Content-Type")
	if ct == "" {
		t.Log("no content-type set; ok")
	}
}
