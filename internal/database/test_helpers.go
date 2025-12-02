package database

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func getTestDB(t *testing.T) *DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set - run 'make test/db'")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	tx := db.MustBegin()
	t.Cleanup(func() {
		tx.Rollback()
		db.Close()
	})

	// Wrap the transaction so each test is isolated
	return &DB{dsn: dsn, conn: tx}
}
