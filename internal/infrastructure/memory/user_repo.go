package memory

import (
	"context"
	"sync"
	"go-auth/internal/domain"
)

// UserRepository is an in-memory implementation of domain.UserRepository.
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User // key: email
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Email] = user
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.users[email]; ok {
		return u, nil
	}
	return nil, nil
}
