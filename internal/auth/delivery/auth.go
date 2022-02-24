package delivery

import (
	"net/http"

	"github.com/crxfoz/otus-hl-network/internal/auth"
	"github.com/crxfoz/otus-hl-network/internal/domain"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authManager auth.AuthManager
	repo        domain.UserUsecase
}

func New(authManager auth.AuthManager, repo domain.UserUsecase) *AuthHandler {
	return &AuthHandler{authManager: authManager, repo: repo}
}

func (ah *AuthHandler) Authorize(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	account, err := ah.repo.FindAccount(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusNotFound, domain.HTTPError{Error: "incorrect username or password"})
	}

	userContext, err := ah.authManager.Generate(domain.User{
		ID:       account.ID,
		Username: account.Username,
		Password: account.Password,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, userContext)
}
