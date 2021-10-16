package domain

import (
	"otus-hl-network/internal/user/repository"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	UserID    int64    `json:"user_id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int64    `json:"age"`
	City      string   `json:"city"`
	Interests []string `json:"interests"`
	Gender    Gender   `json:"gender"`
}

type Friends struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserContext struct {
	ID       int64
	Username string
	Token    string
}

type UserRepo interface {
	FindAccount(username string) (repository.User, error)
	UserListExcept(id ...int64) ([]repository.UserInfo, error)
	FindUserInfo(id int64) (repository.UserInfo, error)
	FindFriends(id int64) ([]repository.UserInfo, error)
	AddUserWithInfo(username string, password string, info repository.UpdateUserInfo) error
	UpdateUserInfo(id int64, info repository.UpdateUserInfo) error
	AddFriends(id int64, friendIDs ...int64) error
}

type UserUsecase interface {
	UserListExcept(id ...int64) ([]UserInfo, error)
}
