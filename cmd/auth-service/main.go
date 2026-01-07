package main

import (
	"context"
	"log/slog"
	"os"

	"time"

	"github.com/gin-gonic/gin"

	"go-auth/internal/app"
	"go-auth/internal/app/usecase"
	"go-auth/internal/config"

	// "go-auth/internal/infrastructure/memory" // Deprecated
	"go-auth/internal/infrastructure/postgres"
	"go-auth/internal/security/jwt"
	"go-auth/internal/security/password"
	httpv1 "go-auth/internal/transport/http"
)

func main() {
	preLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	cfg, err := config.Load()
	if err != nil {
		preLogger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// 2. Setup Logger
	logger := setupLogger(cfg.App.Environment)
	logger.Info("starting auth-service",
		"app", cfg.App.Name,
		"env", cfg.App.Environment,
	)

	// 3. Init Infrastructure
	// DB Pool
	dbPool, err := postgres.InitPool(context.Background(), cfg.Postgres.DSN, logger)
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Migrations выполняются через docker-entrypoint-initdb.d и CI/CD, приложение не исполняет DDL

	userRepo := postgres.NewUserRepository(dbPool)
	refreshRepo := postgres.NewRefreshRepository(dbPool)
	var pwdService app.PasswordService
	if cfg.Security.BcryptCost > 0 {
		pwdService = password.NewWithCost(cfg.Security.BcryptCost)
	} else {
		pwdService = password.New()
	}

	// 4. Init Application / UseCases
	registerUC := usecase.NewRegisterUserUseCase(logger, userRepo, pwdService)

	// Token service and Login use case
	tokenCfg := app.TokenConfig{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
		Issuer:        cfg.App.Name,
		Audience:      cfg.App.Name,
	}
	tokenService := jwt.NewJWTService(tokenCfg)
	loginUC := usecase.NewLoginUserUseCase(logger, userRepo, pwdService, tokenService, refreshRepo)
	refreshUC := usecase.NewRefreshUseCase(tokenService, refreshRepo)
	logoutUC := usecase.NewLogoutUseCase(refreshRepo)

	// 5. Init Transport (HTTP - Gin)
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	r.Use(httpv1.RequestID())
	r.Use(httpv1.CORS())
	r.Use(httpv1.SecurityHeaders())
	r.Use(httpv1.RateLimit(100, time.Minute))

	// Health endpoint at root path for container healthcheck
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// API V1 Group
	v1 := r.Group("/api/v1")

	authHandler := httpv1.NewAuthHandler(logger, registerUC, loginUC, refreshUC, logoutUC)
	authHandler.RegisterRoutes(v1)

	logger.Info("server started", "port", cfg.HTTP.Port)
	if err := r.Run(":" + cfg.HTTP.Port); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	return slog.New(handler)
}
