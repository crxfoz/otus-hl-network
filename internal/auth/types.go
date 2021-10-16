package auth

import (
	"otus-hl-network/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthManager interface {
	Generate(user domain.User) (domain.UserContext, error)
	Verify(accessToken string) (domain.UserContext, error)
}

func EncryptPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
