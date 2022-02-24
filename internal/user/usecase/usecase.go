package usecase

import (
	"fmt"
	"strings"

	"github.com/crxfoz/otus-hl-network/internal/auth"
	"github.com/crxfoz/otus-hl-network/internal/domain"
	"github.com/crxfoz/otus-hl-network/internal/user/repository"
)

type Usecase struct {
	repo domain.UserRepo
}

func (u *Usecase) Search(firstName string, lastName string) ([]domain.UserInfo, error) {
	users, err := u.repo.Search(firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("repository error: %v", err)
	}

	dusers := make([]domain.UserInfo, 0, len(users))
	for _, item := range users {
		dusers = append(dusers, DTOUserInfo(item))
	}

	return dusers, nil
}

func NewUsecase(repo domain.UserRepo) *Usecase {
	return &Usecase{repo: repo}
}

func DTOUserInfo(info repository.UserInfo) domain.UserInfo {
	var interests []string

	if info.Interests.String != "" {
		interests = strings.Split(info.Interests.String, ";")
	}

	return domain.UserInfo{
		UserID: info.UserID,
		UserData: domain.UserData{
			FirstName: info.FirstName.String,
			LastName:  info.LastName.String,
			Age:       info.Age.Int64,
			City:      info.City.String,
			Interests: interests,
			Gender:    domain.Gender(info.Gender.String),
		},
	}
}

func (u *Usecase) UserListExcept(id ...int64) ([]domain.UserInfo, error) {
	users, err := u.repo.UserListExcept(id...)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	out := make([]domain.UserInfo, 0, len(users))
	for i := range users {
		out = append(out, DTOUserInfo(users[i]))
	}

	return out, nil
}

func (u *Usecase) FindFriends(id int64) ([]domain.UserInfo, error) {
	friends, err := u.repo.FindFriends(id)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	out := make([]domain.UserInfo, 0, len(friends))
	for i := range friends {
		out = append(out, DTOUserInfo(friends[i]))
	}

	return out, nil
}

func (u *Usecase) AddFriend(userID int64, friendID int64) error {
	if err := u.repo.AddFriends(userID, friendID); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) DeleteFriend(userID int64, friendID int64) error {
	if err := u.repo.DeleteFriends(userID, friendID); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) FindUserInfo(id int64) (domain.UserInfo, error) {
	profile, err := u.repo.FindUserInfo(id)
	if err != nil {
		return domain.UserInfo{}, err
	}

	return DTOUserInfo(profile), nil
}

func (u *Usecase) FindAccount(username string) (domain.User, error) {
	account, err := u.repo.FindAccount(username)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:       account.ID,
		Username: account.Username,
		Password: account.Password,
	}, nil
}

func (u *Usecase) AddUserWithInfo(username string, password string, info domain.UserData) error {
	var interests string

	if len(info.Interests) > 0 {
		interests = strings.Join(info.Interests, ";")
	}

	return u.repo.AddUserWithInfo(username, auth.EncryptPassword(password), repository.UpdateUserInfo{
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Age:       info.Age,
		City:      info.City,
		Interests: interests,
		Gender:    string(info.Gender),
	})
}

func (u *Usecase) UpdateUserInfo(id int64, info domain.UserData) error {
	var interests string

	if len(info.Interests) > 0 {
		interests = strings.Join(info.Interests, ";")
	}

	return u.repo.UpdateUserInfo(id, repository.UpdateUserInfo{
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Age:       info.Age,
		City:      info.City,
		Interests: interests,
		Gender:    string(info.Gender),
	})
}
