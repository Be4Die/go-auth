package domain

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
	RevokeAllByUser(ctx context.Context, userID string) error
}
