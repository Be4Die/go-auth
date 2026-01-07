package domain

import "context"

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
