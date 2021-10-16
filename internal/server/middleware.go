package server

import (
	"otus-hl-network/internal/auth"
	"otus-hl-network/internal/domain"

	"github.com/labstack/echo"
)

type UserDataNext func(echo.Context, domain.UserContext) error

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func (a *AuthMiddleware) Do(next UserDataNext) echo.HandlerFunc {
	return func(c echo.Context) error {
		// token := c.Request().Header.Get("Authorization")
		// token = strings.Replace(token, "Bearer ", "", 1)
		//
		// logrus.WithField("token", token).Info("got token")
		//
		// claims, err := a.jwtManager.Verify(token)
		// if err != nil {
		// 	return c.JSON(http.StatusForbidden, domain.HTTPError{Error: "bad auth"})
		// }
		//
		// userData := domain.UserContext{
		// 	ID:       claims.UserID,
		// 	Username: claims.Username,
		// 	Token:    token,
		// }

		userData := domain.UserContext{
			ID:       1,
			Username: "foo@gmail.com",
			Token:    "",
		}

		return next(c, userData)
	}
}
