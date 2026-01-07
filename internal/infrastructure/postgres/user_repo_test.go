package postgres

import (
	"context"
	"fmt"
	"go-auth/internal/domain"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestUserRepository_CreateAndFind(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set; skipping integration test")
	}

	log := slog.New(slog.NewTextHandler(testWriter{}, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := InitPool(ctx, dsn, log)
	if err != nil {
		t.Fatalf("init pool: %v", err)
	}
	defer pool.Close()

	_, _ = pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	// Ensure table exists (for CI service container)
	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        is_verified BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	repo := NewUserRepository(pool)
	u := domain.NewUser("int@ex.com", "hash")
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	if u.ID == "" {
		t.Fatalf("id not set after create")
	}

	got, err := repo.FindByEmail(ctx, "int@ex.com")
	if err != nil || got == nil || got.Email != "int@ex.com" {
		t.Fatalf("find: got=%v err=%v", got, err)
	}
	fmt.Println("created user id:", got.ID)
}

type testWriter struct{}

func (testWriter) Write(p []byte) (int, error) { return len(p), nil }
