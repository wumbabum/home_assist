package main

import (
	"log/slog"
	"os"
	"testing"

	"github.com/alexedwards/scs/v2"
)

// newTestApplication creates a minimal application instance for testing
func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	return &application{
		logger: logger,
	}
}

// newTestApplicationWithSession creates application with in-memory session manager
func newTestApplicationWithSession(t *testing.T) *application {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	// Create in-memory session manager for tests
	sessionManager := scs.New()

	return &application{
		logger:         logger,
		sessionManager: sessionManager,
	}
}
