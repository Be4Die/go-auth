package app

import "time"

// TokenService defines the interface for token generation and validation.
type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (string, error) // Returns userID if valid
	ValidateRefresh(token string) (string, error)
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

type TokenConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	Issuer        string
	Audience      string
}
