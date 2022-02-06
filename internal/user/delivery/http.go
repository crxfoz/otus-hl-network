package delivery

import (
	"net/http"
	"strconv"

	"otus-hl-network/internal/domain"

	"github.com/labstack/echo"
)

type UserHander struct {
	repo domain.UserUsecase
}

func New(ua domain.UserUsecase) *UserHander {
	return &UserHander{repo: ua}
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	domain.UserData
}

func (h *UserHander) Register(c echo.Context) error {
	var userReg UserRegisterRequest

	if err := c.Bind(&userReg); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.HTTPError{Error: "incorect body"})
	}

	if err := h.repo.AddUserWithInfo(userReg.Username, userReg.Password, userReg.UserData); err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	return nil
}

func (h *UserHander) UpdateProfile(c echo.Context, userContext domain.UserContext) error {
	var userData domain.UserData

	if err := c.Bind(&userData); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.HTTPError{Error: "incorect body"})
	}

	if err := h.repo.UpdateUserInfo(userContext.ID, userData); err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	return nil
}

func (h *UserHander) Profile(c echo.Context, userContext domain.UserContext) error {
	var profileID int64

	id := c.Param("id")
	if id == "" {
		profileID = userContext.ID
	} else {
		idParsed, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.JSON(http.StatusNotFound, domain.HTTPError{Error: "used id isn't valide"})
		}

		profileID = idParsed
	}

	profile, err := h.repo.FindUserInfo(profileID)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.HTTPError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *UserHander) Users(c echo.Context, userContext domain.UserContext) error {
	users, err := h.repo.UserListExcept(userContext.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHander) Friends(c echo.Context, userContext domain.UserContext) error {
	friends, err := h.repo.FindFriends(userContext.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, friends)
}

func (h *UserHander) AddFriend(c echo.Context, userContext domain.UserContext) error {
	id := c.Param("id")
	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.HTTPError{Error: "used id isn't valide"})
	}

	if err := h.repo.AddFriend(userContext.ID, idParsed); err != nil {
		return c.JSON(http.StatusNotFound, domain.HTTPError{Error: "user with given id doesn't exist"})
	}

	return c.JSON(http.StatusOK, domain.HTTPok{Status: "ok"})
}

func (h *UserHander) DeleteFriend(c echo.Context, userContext domain.UserContext) error {
	id := c.Param("id")
	idParsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.HTTPError{Error: "used id isn't valide"})
	}

	if err := h.repo.DeleteFriend(userContext.ID, idParsed); err != nil {
		return c.JSON(http.StatusNotFound, domain.HTTPError{Error: "user with given id doesn't exist"})
	}

	return c.JSON(http.StatusOK, domain.HTTPok{Status: "ok"})
}

func (h *UserHander) Search(c echo.Context, userContext domain.UserContext) error {
	firstName := c.QueryParam("fname")
	lastName := c.QueryParam("lname")

	users, err := h.repo.Search(firstName, lastName)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.HTTPError{Error: "users with given params aren't found"})

	}

	return c.JSON(http.StatusOK, users)
}
