package jwt

import (
	"go-auth/internal/app"
	"testing"
	"time"
)

func TestGenerateAndValidateTokens(t *testing.T) {
	cfg := app.TokenConfig{
		AccessSecret:  "test-access",
		RefreshSecret: "test-refresh",
		AccessTTL:     1 * time.Minute,
		RefreshTTL:    10 * time.Minute,
	}

	s := NewJWTService(cfg)

	access, err := s.GenerateAccessToken("user-1")
	if err != nil || access == "" {
		t.Fatalf("failed to generate access token: %v", err)
	}

	uid, err := s.ValidateToken(access)
	if err != nil || uid != "user-1" {
		t.Fatalf("validate failed, uid=%s err=%v", uid, err)
	}

	refresh, err := s.GenerateRefreshToken("user-1")
	if err != nil || refresh == "" {
		t.Fatalf("failed to generate refresh token: %v", err)
	}
}
