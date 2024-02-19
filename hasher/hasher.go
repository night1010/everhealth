package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(text string) ([]byte, error)
	Compare(hashed string, text string) bool
}

type bcryptHasher struct{}

func NewHasher() Hasher {
	return &bcryptHasher{}
}

func (h *bcryptHasher) Hash(text string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
}

func (h *bcryptHasher) Compare(hashed string, text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(text))
	return err == nil
}
