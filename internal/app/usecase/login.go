package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go-auth/internal/app"
	"go-auth/internal/domain"
	"go-auth/internal/security/tokenhash"
)

type LoginUserCmd struct {
	Email    string
	Password string
}

type LoginUserResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type LoginUserUseCase struct {
	log          *slog.Logger
	userRepo     domain.UserRepository
	pwdService   app.PasswordService
	tokenService app.TokenService
	refreshRepo  domain.RefreshTokenRepository
}

func NewLoginUserUseCase(
	log *slog.Logger,
	userRepo domain.UserRepository,
	pwdService app.PasswordService,
	tokenService app.TokenService,
	refreshRepo domain.RefreshTokenRepository,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		log:          log,
		userRepo:     userRepo,
		pwdService:   pwdService,
		tokenService: tokenService,
		refreshRepo:  refreshRepo,
	}
}

func (uc *LoginUserUseCase) Handle(ctx context.Context, cmd LoginUserCmd) (*LoginUserResult, error) {
	log := uc.log.With("op", "LoginUser", "email", cmd.Email)

	// 1. Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, app.NewError(app.ErrCodeInvalidCredentials, "Invalid credentials")
	}

	// 2. Verify password
	err = uc.pwdService.Compare(user.Password, cmd.Password)
	if err != nil {
		log.Warn("invalid password attempt")
		return nil, app.NewError(app.ErrCodeInvalidCredentials, "Invalid credentials")
	}

	// 3. Generate tokens
	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if uc.refreshRepo != nil {
		_ = uc.refreshRepo.Save(ctx, user.ID, tokenhash.Hash(refreshToken), time.Now().Add(uc.tokenService.RefreshTTL()))
	}

	log.Info("user logged in successfully", "user_id", user.ID)

	return &LoginUserResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(uc.tokenService.AccessTTL().Seconds()),
	}, nil
}

func (uc *LoginUserUseCase) TokenUserID(token string) (string, error) {
	return uc.tokenService.ValidateToken(token)
}
