package main

import (
	"log/slog"
	"os"
	"testing"
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
