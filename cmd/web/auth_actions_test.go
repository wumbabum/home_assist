package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wumbabum/home_assist/internal/authenticator"
)

func TestLogin(t *testing.T) {
	app := &application{
		auth0: &authenticator.Authenticator{},
	}

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	app.login(w, req)

	// Just verify redirect happens
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, w.Code)
	}
}

func TestCallback(t *testing.T) {
	app := &application{}
	req := httptest.NewRequest(http.MethodGet, "/callback", nil)
	w := httptest.NewRecorder()

	app.callback(w, req)

	// Just verify stub redirects home
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
	}
}

func TestLogout(t *testing.T) {
	app := &application{}
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	app.logout(w, req)

	// Just verify stub redirects home
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
	}
}
