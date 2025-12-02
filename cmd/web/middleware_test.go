package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuth_NoSession(t *testing.T) {
	// Test that requireAuth redirects to /login when no session exists
	app := newTestApplicationWithSession(t)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called when not authenticated")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Wrap with session middleware first, then auth
	handler := app.sessionManager.LoadAndSave(app.requireAuth(testHandler))
	handler.ServeHTTP(w, req)

	// Should redirect to login
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	if location := w.Header().Get("Location"); location != "/login" {
		t.Errorf("expected redirect to /login, got %s", location)
	}
}

func TestRecoverPanic(t *testing.T) {
	app := newTestApplication(t)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	wrappedHandler := app.recoverPanic(panicHandler)
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestSecurityHeaders(t *testing.T) {
	app := &application{}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	wrappedHandler := app.securityHeaders(testHandler)
	wrappedHandler.ServeHTTP(w, req)

	tests := []struct {
		header   string
		expected string
	}{
		{"Referrer-Policy", "origin-when-cross-origin"},
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "deny"},
	}

	for _, tt := range tests {
		got := w.Header().Get(tt.header)
		if got != tt.expected {
			t.Errorf("header %s: expected %q, got %q", tt.header, tt.expected, got)
		}
	}
}

func TestLogAccess(t *testing.T) {
	app := newTestApplication(t)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	wrappedHandler := app.logAccess(testHandler)
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
