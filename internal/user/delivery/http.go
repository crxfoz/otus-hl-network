package delivery

import (
	"net/http"

	"otus-hl-network/internal/domain"

	"github.com/labstack/echo"
)

type UserHander struct {
	repo domain.UserUsecase
}

func New(ua domain.UserUsecase) *UserHander {
	return &UserHander{repo: ua}
}

// func (h *UserHander) Authorize(c echo.Context) error {
// }
//
// func (h *UserHander) Login(c echo.Context) error {
// }
//
// func (h *UserHander) Register(c echo.Context) error {
//
// }
//
// func (h *UserHander) Profile(c echo.Context, userContext domain.UserContext) error {
//
// }

func (h *UserHander) Users(c echo.Context, userContext domain.UserContext) error {
	users, err := h.repo.UserListExcept(userContext.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// func (h *UserHander) Friends(c echo.Context, userContext domain.UserContext) error {
//
// }
//
// func (h *UserHander) AddFriend(c echo.Context, userContext domain.UserContext) error {
//
// }
//
// func (h *UserHander) DeleteFriend(c echo.Context, userContext domain.UserContext) error {
//
// }
