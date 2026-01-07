package httpv1

import (
	"log/slog"
	"net/http"

	"go-auth/internal/app"
	"go-auth/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	log        *slog.Logger
	registerUC *usecase.RegisterUserUseCase
	loginUC    *usecase.LoginUserUseCase
	refreshUC  *usecase.RefreshUseCase
	logoutUC   *usecase.LogoutUseCase
}

func NewAuthHandler(
	log *slog.Logger,
	registerUC *usecase.RegisterUserUseCase,
	loginUC *usecase.LoginUserUseCase,
	refreshUC *usecase.RefreshUseCase,
	logoutUC *usecase.LogoutUseCase,
) *AuthHandler {
	return &AuthHandler{
		log:        log,
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
		logoutUC:   logoutUC,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
	}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := usecase.RegisterUserCmd{
		Email:    req.Email,
		Password: req.Password,
	}
	if !validPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not meet complexity", "code": app.ErrCodeValidation})
		return
	}

	if err := h.registerUC.Handle(c.Request.Context(), cmd); err != nil {
		h.log.Error("registration failed", "error", err)
		code := app.ErrCodeInternal
		status := http.StatusInternalServerError
		msg := "Registration failed"
		if ae, ok := err.(app.AppError); ok {
			code = ae.Code
			msg = ae.Msg
			switch ae.Code {
			case app.ErrCodeEmailExists:
				status = http.StatusConflict
			case app.ErrCodeValidation:
				status = http.StatusBadRequest
			}
		}
		c.JSON(status, gin.H{"error": msg, "code": code})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := usecase.LoginUserCmd{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.loginUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.log.Warn("login failed", "error", err)
		status := http.StatusUnauthorized
		code := app.ErrCodeInvalidCredentials
		msg := "Invalid credentials"
		if ae, ok := err.(app.AppError); ok {
			code = ae.Code
			msg = ae.Msg
		}
		c.JSON(status, gin.H{"error": msg, "code": code})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  res.AccessToken,
		"refresh_token": res.RefreshToken,
		"expires_in":    res.ExpiresIn,
		"token_type":    "Bearer",
	})
}

func validPassword(p string) bool {
	var up, low, dig, spec bool
	for i := 0; i < len(p); i++ {
		c := p[i]
		if c >= 'A' && c <= 'Z' {
			up = true
		}
		if c >= 'a' && c <= 'z' {
			low = true
		}
		if c >= '0' && c <= '9' {
			dig = true
		}
		if (c < '0' || c > '9') && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
			spec = true
		}
	}
	return len(p) >= 8 && up && low && dig && spec
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": app.ErrCodeValidation})
		return
	}
	res, err := h.refreshUC.Handle(c.Request.Context(), usecase.RefreshCmd{RefreshToken: req.RefreshToken})
	if err != nil {
		h.log.Warn("refresh failed", "error", err)
		status := http.StatusUnauthorized
		code := app.ErrCodeInvalidCredentials
		msg := "Invalid refresh token"
		if ae, ok := err.(app.AppError); ok {
			code = ae.Code
			msg = ae.Msg
		}
		c.JSON(status, gin.H{"error": msg, "code": code})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": res.AccessToken, "refresh_token": res.RefreshToken, "expires_in": res.ExpiresIn, "token_type": "Bearer"})
}

func (h *AuthHandler) logout(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 || auth[:7] != "Bearer " {
		c.Status(http.StatusUnauthorized)
		return
	}
	token := auth[7:]
	uid, err := h.loginUC.TokenUserID(token)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	if err := h.logoutUC.Handle(c.Request.Context(), usecase.LogoutCmd{UserID: uid}); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}
