package usecase

import (
	"context"
	"errors"
	"go-auth/internal/domain"
	"log/slog"
	"testing"
)

type fakePwd struct{}

func (f *fakePwd) Hash(p string) (string, error) { return "hash:" + p, nil }
func (f *fakePwd) Compare(h, p string) error {
	if h == "hash:"+p {
		return nil
	}
	return errors.New("mismatch")
}

type memRepo struct{ users map[string]*domain.User }

func (r *memRepo) Create(ctx context.Context, u *domain.User) error {
	if r.users == nil {
		r.users = map[string]*domain.User{}
	}
	r.users[u.Email] = u
	u.ID = "id-1"
	return nil
}
func (r *memRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if r.users == nil {
		return nil, nil
	}
	return r.users[email], nil
}

func TestRegister_Success(t *testing.T) {
	log := slog.New(slog.NewTextHandler(testWriter{}, nil))
	uc := NewRegisterUserUseCase(log, &memRepo{}, &fakePwd{})
	err := uc.Handle(context.Background(), RegisterUserCmd{Email: "u@ex.com", Password: "p"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
}

func TestRegister_Duplicate(t *testing.T) {
	log := slog.New(slog.NewTextHandler(testWriter{}, nil))
	repo := &memRepo{}
	uc := NewRegisterUserUseCase(log, repo, &fakePwd{})
	_ = uc.Handle(context.Background(), RegisterUserCmd{Email: "u@ex.com", Password: "p"})
	if err := uc.Handle(context.Background(), RegisterUserCmd{Email: "u@ex.com", Password: "p"}); err == nil {
		t.Fatalf("expected duplicate error")
	}
}

// testWriter discards logs
type testWriter struct{}

func (testWriter) Write(p []byte) (int, error) { return len(p), nil }
