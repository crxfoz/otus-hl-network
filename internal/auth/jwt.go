package auth

import (
	"fmt"
	"time"

	"github.com/crxfoz/otus-hl-network/internal/domain"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}

type JWTManager struct {
	tokenDuration time.Duration
	secretKey     string
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (manager *JWTManager) Generate(user domain.User) (domain.UserContext, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Username,
		UserID:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenKey, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		return domain.UserContext{}, err
	}

	return domain.UserContext{
		ID:       user.ID,
		Username: user.Username,
		Token:    tokenKey,
	}, nil
}

func (manager *JWTManager) Verify(accessToken string) (domain.UserContext, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return domain.UserContext{}, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return domain.UserContext{}, fmt.Errorf("invalid token claims")
	}

	return domain.UserContext{
		ID:       claims.UserID,
		Username: claims.Username,
		Token:    accessToken,
	}, nil
}
