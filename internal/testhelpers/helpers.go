package testhelpers

import (
	"os"
	"testing"

	"github.com/wumbabum/home_assist/internal/database"
)

func GetTestDB(t *testing.T) *database.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set - run 'make test/db'")
	}

	db, err := database.New(dsn)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
