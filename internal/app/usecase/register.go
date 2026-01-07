package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"go-auth/internal/app"
	"go-auth/internal/domain"
)

type RegisterUserCmd struct {
	Email    string
	Password string
}

type RegisterUserUseCase struct {
	log        *slog.Logger
	userRepo   domain.UserRepository
	pwdService app.PasswordService
}

func NewRegisterUserUseCase(log *slog.Logger, userRepo domain.UserRepository, pwdService app.PasswordService) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		log:        log,
		userRepo:   userRepo,
		pwdService: pwdService,
	}
}

func (uc *RegisterUserUseCase) Handle(ctx context.Context, cmd RegisterUserCmd) error {
	log := uc.log.With("op", "RegisterUser", "email", cmd.Email)

	// 1. Check if user exists
	// Note: In a real DB impl, FindByEmail might return a specific "NotFound" error.
	// Here we assume if err == nil and user != nil, then user exists.
	existing, err := uc.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		// If error is NOT "not found", return it
		// For now, let's assume we handle this in the repo implementation or specific error types
		// This is a simplification.
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if existing != nil {
		return app.NewError(app.ErrCodeEmailExists, "Email already exists")
	}

	// 2. Hash password
	hash, err := uc.pwdService.Hash(cmd.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.NewUser(cmd.Email, hash)

	// 4. Save to repo
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Info("user registered successfully", "user_id", user.ID)
	return nil
}
