package hasher

import (
	"errors"
	"fmt"
	"github.com/BernsteinMondy/medods-test-task/src/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type HasherService struct {
}

var _ service.Hasher = (*HasherService)(nil)

func NewService() *HasherService {
	return &HasherService{}
}

func (h *HasherService) Hash(src string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash: %w", err)
	}

	return string(hash), nil
}

func (h *HasherService) CompareHashAndPassword(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("compare hash and password: %w", err)
	}

	return true, nil
}
