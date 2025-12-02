package database

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Auth0Sub  string    `db:"auth0_sub"` // Auth0 subject ID in auth0.dev
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Picture   string    `db:"picture"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (db *DB) UpsertUser(ctx context.Context, auth0Sub, email, name, picture string) (*User, error) {
	query := `
		INSERT INTO users (auth0_sub, email, name, picture, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (auth0_sub)
		DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			picture = EXCLUDED.picture,
			updated_at = NOW()
		RETURNING id, auth0_sub, email, name, picture, created_at, updated_at
	`
	var user User
	err := db.Get(&user, query, auth0Sub, email, name, picture)
	return &user, err
}

func (db *DB) GetUserBySub(ctx context.Context, auth0Sub string) (*User, error) {
	query := `SELECT id, auth0_sub, email, name, picture, created_at, updated_at FROM users WHERE auth0_sub = $1`
	var user User
	err := db.Get(&user, query, auth0Sub)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
