package server

import (
	"net/http"
	"strings"

	"github.com/crxfoz/otus-hl-network/internal/auth"
	"github.com/crxfoz/otus-hl-network/internal/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type UserDataNext func(echo.Context, domain.UserContext) error

type AuthMiddleware struct {
	authManager auth.AuthManager
}

func (a *AuthMiddleware) Do(next UserDataNext) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)

		logrus.WithField("token", token).Info("got token")

		claims, err := a.authManager.Verify(token)
		if err != nil {
			return c.JSON(http.StatusForbidden, domain.HTTPError{Error: "bad auth"})
		}

		userData := domain.UserContext{
			ID:       claims.ID,
			Username: claims.Username,
			Token:    token,
		}

		return next(c, userData)
	}
}
