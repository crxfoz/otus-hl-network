package repository

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// TODO: use auto-integration tests and unit tests with sqlmock

func TestUserRepo_FindAccount(t *testing.T) {
	db, err := sqlx.Connect("mysql", "root:user123@(localhost:3306)/hl_network")
	assert.Nil(t, err)

	repo := NewUserRepo(db)

	user, err := repo.FindAccount("foo@gmail.com")
	assert.Nil(t, err)

	fmt.Println(user)
}

func TestUserRepo_FindUserInfo(t *testing.T) {
	db, err := sqlx.Connect("mysql", "root:user123@(localhost:3306)/hl_network")
	assert.Nil(t, err)

	repo := NewUserRepo(db)

	user, err := repo.FindUserInfo(1)
	assert.Nil(t, err)

	fmt.Println(user)
}

func TestUserRepo_FindFriends(t *testing.T) {
	db, err := sqlx.Connect("mysql", "root:user123@(localhost:3306)/hl_network")
	assert.Nil(t, err)

	repo := NewUserRepo(db)

	user, err := repo.FindFriends(1)
	assert.Nil(t, err)

	fmt.Println(user)
}

func TestUserRepo_UpdateUserInfo(t *testing.T) {
	db, err := sqlx.Connect("mysql", "root:user123@(localhost:3306)/hl_network")
	assert.Nil(t, err)

	repo := NewUserRepo(db)

	err = repo.UpdateUserInfo(1, UpdateUserInfo{
		FirstName: "",
		LastName:  "",
		Age:       0,
		City:      "",
		Interests: "",
		Gender:    "",
	})
	assert.Nil(t, err)
}

func TestUserRepo_AddUserWithInfo(t *testing.T) {
	db, err := sqlx.Connect("mysql", "root:user123@(localhost:3306)/hl_network")
	assert.Nil(t, err)

	repo := NewUserRepo(db)

	err = repo.AddUserWithInfo("kek@sobaka.kosha", "awesome", UpdateUserInfo{
		FirstName: "Koshka",
		LastName:  "Cat",
		Age:       0,
		City:      "",
		Interests: "",
		Gender:    "",
	})
	assert.Nil(t, err)
}

func TestUserRepo_AddFriends(t *testing.T) {
	db, err := sqlx.Connect("mysql", "root:user123@(localhost:3306)/hl_network")
	assert.Nil(t, err)

	repo := NewUserRepo(db)

	err = repo.AddFriends(5, 1, 3)
	assert.Nil(t, err)
}
