package usecase

import (
	"context"
	"go-auth/internal/app"
	"go-auth/internal/domain"
	"go-auth/internal/security/tokenhash"
	"time"
)

type RefreshCmd struct{ RefreshToken string }

type RefreshUseCase struct {
	tokens app.TokenService
	repo   domain.RefreshTokenRepository
}

func NewRefreshUseCase(tokens app.TokenService, repo domain.RefreshTokenRepository) *RefreshUseCase {
	return &RefreshUseCase{tokens: tokens, repo: repo}
}

func (uc *RefreshUseCase) Handle(ctx context.Context, cmd RefreshCmd) (*LoginUserResult, error) {
	uid, err := uc.tokens.ValidateRefresh(cmd.RefreshToken)
	if err != nil {
		return nil, app.NewError(app.ErrCodeInvalidCredentials, "Invalid refresh token")
	}
	h := tokenhash.Hash(cmd.RefreshToken)
	rec, err := uc.repo.FindByHash(ctx, h)
	if rec == nil {
		return nil, app.NewError(app.ErrCodeInvalidCredentials, "Invalid refresh token")
	}
	if rec.RevokedAt != nil || time.Now().After(rec.ExpiresAt) {
		return nil, app.NewError(app.ErrCodeInvalidCredentials, "Invalid refresh token")
	}

	_ = uc.repo.RevokeByHash(ctx, h)

	access, err := uc.tokens.GenerateAccessToken(uid)
	if err != nil {
		return nil, app.NewError(app.ErrCodeInternal, "Failed to generate access token")
	}
	newRefresh, err := uc.tokens.GenerateRefreshToken(uid)
	if err != nil {
		return nil, app.NewError(app.ErrCodeInternal, "Failed to generate refresh token")
	}
	_ = uc.repo.Save(ctx, uid, tokenhash.Hash(newRefresh), time.Now().Add(24*time.Hour*7))

	return &LoginUserResult{AccessToken: access, RefreshToken: newRefresh, ExpiresIn: int64(uc.tokens.AccessTTL().Seconds())}, nil
}

type LogoutCmd struct{ UserID string }

type LogoutUseCase struct{ repo domain.RefreshTokenRepository }

func NewLogoutUseCase(repo domain.RefreshTokenRepository) *LogoutUseCase {
	return &LogoutUseCase{repo: repo}
}

func (uc *LogoutUseCase) Handle(ctx context.Context, cmd LogoutCmd) error {
	return uc.repo.RevokeAllByUser(ctx, cmd.UserID)
}
