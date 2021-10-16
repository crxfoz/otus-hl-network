package usecase

import (
	"fmt"
	"strings"

	"otus-hl-network/internal/domain"
)

type Usecase struct {
	repo domain.UserRepo
}

func NewUsecase(repo domain.UserRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) UserListExcept(id ...int64) ([]domain.UserInfo, error) {
	users, err := u.repo.UserListExcept(id...)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	out := make([]domain.UserInfo, 0, len(users))

	for i := range users {
		var interests []string

		if users[i].Interests.String != "" {
			interests = strings.Split(users[i].Interests.String, ";")
		}

		out = append(out, domain.UserInfo{
			UserID:    users[i].UserID,
			FirstName: users[i].FirstName.String,
			LastName:  users[i].LastName.String,
			Age:       users[i].Age.Int64,
			City:      users[i].City.String,
			Interests: interests,
			Gender:    domain.Gender(users[i].Gender.String),
		})
	}

	return out, nil
}
