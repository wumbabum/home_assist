package database

import (
	"os"
	"testing"
)

func getTestDB(t *testing.T) *DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set - run 'make test/db'")
	}

	db, err := New(dsn)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
