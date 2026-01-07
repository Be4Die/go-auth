package jwt

import (
	"fmt"
	"time"

	"go-auth/internal/app"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	config app.TokenConfig
}

func NewJWTService(cfg app.TokenConfig) *JWTService {
	return &JWTService{
		config: cfg,
	}
}

func (s *JWTService) GenerateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, s.config.AccessSecret, s.config.AccessTTL)
}

func (s *JWTService) GenerateRefreshToken(userID string) (string, error) {
	return s.generateToken(userID, s.config.RefreshSecret, s.config.RefreshTTL)
}

func (s *JWTService) generateToken(userID string, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
		"iss": s.config.Issuer,
		"aud": s.config.Audience,
		"jti": fmt.Sprintf("%s-%d", userID, time.Now().UnixNano()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *JWTService) ValidateToken(tokenString string) (string, error) {
	// Note: This needs to know which secret to use.
	// For simplicity in this example, we'll assume access token validation primarily.
	// In a real scenario, you might have separate methods or pass the secret type.
	return s.validate(tokenString, s.config.AccessSecret)
}

func (s *JWTService) ValidateRefresh(tokenString string) (string, error) {
	return s.validate(tokenString, s.config.RefreshSecret)
}

func (s *JWTService) validate(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", fmt.Errorf("invalid token claims")
}

func (s *JWTService) AccessTTL() time.Duration  { return s.config.AccessTTL }
func (s *JWTService) RefreshTTL() time.Duration { return s.config.RefreshTTL }
