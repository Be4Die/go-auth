package httpv1

import (
    "context"
    "fmt"
    "log/slog"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "go-auth/internal/app"
    "go-auth/internal/app/usecase"
    "go-auth/internal/domain"

    "github.com/gin-gonic/gin"
)

type memRepo struct{ users map[string]*domain.User }

func (r *memRepo) Create(_ context.Context, u *domain.User) error {
	if r.users == nil {
		r.users = map[string]*domain.User{}
	}
	r.users[u.Email] = u
	u.ID = "id-1"
	return nil
}
func (r *memRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	if r.users == nil {
		return nil, nil
	}
	return r.users[email], nil
}

type fakePwd struct{}

func (fakePwd) Hash(p string) (string, error) { return "hash:" + p, nil }
func (fakePwd) Compare(h, p string) error {
	if h == "hash:"+p {
		return nil
	}
	return fmt.Errorf("bad")
}

type fakeToken struct{}

func (fakeToken) GenerateAccessToken(userID string) (string, error)  { return "acc:" + userID, nil }
func (fakeToken) GenerateRefreshToken(userID string) (string, error) { return "ref:" + userID, nil }
func (fakeToken) ValidateToken(token string) (string, error)         { return "", nil }
func (fakeToken) ValidateRefresh(token string) (string, error)       { return "id-1", nil }
func (fakeToken) AccessTTL() time.Duration                           { return time.Minute }
func (fakeToken) RefreshTTL() time.Duration                          { return 7 * 24 * time.Hour }

func TestRoutes_RegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &memRepo{}
	regUC := usecase.NewRegisterUserUseCase(slog.Default(), repo, app.PasswordService(fakePwd{}))
    logUC := usecase.NewLoginUserUseCase(slog.Default(), repo, app.PasswordService(fakePwd{}), app.TokenService(fakeToken{}), nil)

    h := NewAuthHandler(slog.Default(), regUC, logUC, nil, nil)
	h.RegisterRoutes(r.Group("/api/v1"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(`{"email":"t@e.com","password":"Password123!"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("register code=%d", w.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"t@e.com","password":"Password123!"}`))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	if w2.Code != 200 {
		t.Fatalf("login code=%d", w2.Code)
	}
}
