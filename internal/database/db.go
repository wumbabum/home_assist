package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wumbabum/home_assist/assets"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	dsn string
	*sqlx.DB
}

func New(dsn string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Add postgres:// prefix if not already present
	connStr := dsn
	if !strings.HasPrefix(dsn, "postgres://") && !strings.HasPrefix(dsn, "postgresql://") {
		connStr = "postgres://" + dsn
	}

	// Debug: print connection string
	fmt.Printf("DEBUG: Connecting with DSN: %s\n", connStr)

	db, err := sqlx.ConnectContext(ctx, "postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	return &DB{dsn: dsn, DB: db}, nil
}

func (db *DB) MigrateUp() error {
	iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
	if err != nil {
		return err
	}

	// Add postgres:// prefix if not already present
	connStr := db.dsn
	if !strings.HasPrefix(db.dsn, "postgres://") && !strings.HasPrefix(db.dsn, "postgresql://") {
		connStr = "postgres://" + db.dsn
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, connStr)
	if err != nil {
		return err
	}

	err = migrator.Up()
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		return nil
	default:
		return err
	}
}
