package usecase

import (
    "context"
    "errors"
    "time"
    "go-auth/internal/app"
    "go-auth/internal/domain"
    "log/slog"
    "testing"
)

type fakeToken struct{}

func (fakeToken) GenerateAccessToken(userID string) (string, error)  { return "acc:" + userID, nil }
func (fakeToken) GenerateRefreshToken(userID string) (string, error) { return "ref:" + userID, nil }
func (fakeToken) ValidateToken(token string) (string, error)         { return "", nil }
func (fakeToken) ValidateRefresh(token string) (string, error)       { return "id-1", nil }
func (fakeToken) AccessTTL() time.Duration                           { return time.Minute }
func (fakeToken) RefreshTTL() time.Duration                          { return 7 * 24 * time.Hour }

type fakePwd2 struct{}

func (fakePwd2) Hash(p string) (string, error) { return "", nil }
func (fakePwd2) Compare(h, p string) error {
	if h == p {
		return nil
	}
	return errors.New("bad")
}

type memRepo2 struct{ u *domain.User }

func (r *memRepo2) Create(ctx context.Context, u *domain.User) error { r.u = u; return nil }
func (r *memRepo2) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.u, nil
}

func TestLogin_Success(t *testing.T) {
	log := slog.New(slog.NewTextHandler(testWriter{}, nil))
	repo := &memRepo2{u: &domain.User{ID: "id-1", Email: "u@ex.com", Password: "p"}}
    uc := NewLoginUserUseCase(log, repo, app.PasswordService(fakePwd2{}), app.TokenService(fakeToken{}), nil)
	res, err := uc.Handle(context.Background(), LoginUserCmd{Email: "u@ex.com", Password: "p"})
	if err != nil || res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatalf("login failed: %v", err)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	log := slog.New(slog.NewTextHandler(testWriter{}, nil))
	repo := &memRepo2{u: &domain.User{ID: "id-1", Email: "u@ex.com", Password: "p"}}
    uc := NewLoginUserUseCase(log, repo, app.PasswordService(fakePwd2{}), app.TokenService(fakeToken{}), nil)
	if _, err := uc.Handle(context.Background(), LoginUserCmd{Email: "u@ex.com", Password: "wrong"}); err == nil {
		t.Fatalf("expected invalid credentials")
	}
}
