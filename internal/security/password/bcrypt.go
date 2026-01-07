package password

import "golang.org/x/crypto/bcrypt"

type BcryptService struct{ cost int }

func New() *BcryptService                 { return &BcryptService{cost: bcrypt.DefaultCost} }
func NewWithCost(cost int) *BcryptService { return &BcryptService{cost: cost} }

func (s *BcryptService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	return string(bytes), err
}

func (s *BcryptService) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
