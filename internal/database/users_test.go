package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

func TestUpsertUser(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()
	auth0Sub := "test|upsert-" + t.Name()

	// Create user
	user1, err := db.UpsertUser(ctx, auth0Sub, "test@example.com", "Test User", "https://pic.jpg")
	if err != nil {
		t.Fatal(err)
	}

	if user1.ID == 0 {
		t.Error("expected ID to be set")
	}
	if user1.Auth0Sub != auth0Sub {
		t.Errorf("expected auth0_sub %s, got %s", auth0Sub, user1.Auth0Sub)
	}
	if user1.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user1.Email)
	}
	if user1.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	time.Sleep(10 * time.Millisecond)

	// Update user
	user2, err := db.UpsertUser(ctx, auth0Sub, "updated@example.com", "Updated User", "")
	if err != nil {
		t.Fatal(err)
	}

	if user2.ID != user1.ID {
		t.Errorf("expected ID %d, got %d", user1.ID, user2.ID)
	}
	if user2.Email != "updated@example.com" {
		t.Errorf("expected email updated@example.com, got %s", user2.Email)
	}
	if !user2.CreatedAt.Equal(user1.CreatedAt) {
		t.Error("expected CreatedAt to remain unchanged")
	}
	if !user2.UpdatedAt.After(user1.UpdatedAt) {
		t.Error("expected UpdatedAt to be after original")
	}
}

func TestGetUserBySub(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()
	auth0Sub := "test|get-" + t.Name()

	// Create user
	created, err := db.UpsertUser(ctx, auth0Sub, "get@example.com", "Get User", "")
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve user
	retrieved, err := db.GetUserBySub(ctx, auth0Sub)
	if err != nil {
		t.Fatal(err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, retrieved.ID)
	}
	if retrieved.Auth0Sub != auth0Sub {
		t.Errorf("expected auth0_sub %s, got %s", auth0Sub, retrieved.Auth0Sub)
	}

	// Test non-existent user
	_, err = db.GetUserBySub(ctx, "nonexistent|sub")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}
