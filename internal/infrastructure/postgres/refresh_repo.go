package postgres

import (
	"context"
	"fmt"
	"go-auth/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshRepository(pool *pgxpool.Pool) *RefreshRepository {
	return &RefreshRepository{pool: pool}
}

func (r *RefreshRepository) Save(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES($1,$2,$3)`, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("postgres: save refresh: %w", err)
	}
	return nil
}

func (r *RefreshRepository) FindByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.pool.QueryRow(ctx, `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at FROM refresh_tokens WHERE token_hash=$1`, tokenHash).Scan(
		&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.RevokedAt, &rt.CreatedAt,
	)
	if err != nil {
		return nil, nil
	}
	return &rt, nil
}

func (r *RefreshRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE token_hash=$1`, tokenHash)
	return err
}

func (r *RefreshRepository) RevokeAllByUser(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE user_id=$1 AND revoked_at IS NULL`, userID)
	return err
}
