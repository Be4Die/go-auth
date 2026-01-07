package app

// PasswordService defines the interface for password hashing and verification.
type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}
