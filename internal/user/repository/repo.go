package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepo struct {
	conn *sqlx.DB
}

func (r *UserRepo) Search(firstName string, lastName string) ([]UserInfo, error) {
	var accounts []UserInfo

	if err := r.conn.Select(&accounts, "SELECT * FROM user_info WHERE first_name LIKE ? AND last_name LIKE ? ORDER BY id",
		fmt.Sprintf("%s%%", firstName),
		fmt.Sprintf("%s%%", lastName),
	); err != nil {
		return []UserInfo{}, err
	}

	return accounts, nil
}

func NewUserRepo(conn *sqlx.DB) *UserRepo {
	return &UserRepo{conn: conn}
}

func (r *UserRepo) FindAccount(username string) (User, error) {
	var user User

	if err := r.conn.Get(&user, "SELECT * FROM users WHERE username=?", username); err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *UserRepo) UserListExcept(id ...int64) ([]UserInfo, error) {
	if len(id) == 0 {
		id = []int64{0}
	}

	var infos []UserInfo

	query, args, err := sqlx.In("SELECT * FROM user_info WHERE user_id NOT IN (?);", id)
	if err != nil {
		return nil, err
	}

	query = r.conn.Rebind(query)

	if err := r.conn.Select(&infos, query, args...); err != nil {
		return nil, err
	}

	return infos, nil
}

func (r *UserRepo) FindUserInfo(id int64) (UserInfo, error) {
	var info UserInfo

	if err := r.conn.Get(&info, "SELECT * FROM user_info WHERE user_id=?", id); err != nil {
		return UserInfo{}, err
	}

	return info, nil
}

func (r *UserRepo) FindFriends(id int64) ([]UserInfo, error) {
	var infos []UserInfo

	if err := r.conn.Select(&infos, "SELECT ui.* FROM user_info ui LEFT JOIN friends ON ui.user_id = friends.friend_id WHERE friends.user_id=?", id); err != nil {
		return nil, err
	}

	return infos, nil
}

func (r *UserRepo) AddUserWithInfo(username string, password string, info UpdateUserInfo) error {
	tx, err := r.conn.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newID int64

	if err := tx.Get(&newID, "SELECT LAST_INSERT_ID()"); err != nil {
		logrus.Info("last intesrt id: ", newID)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO user_info (user_id, first_name, last_name, age, interests, city, gender) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newID,
		info.FirstName,
		info.LastName,
		info.Age,
		info.Interests,
		info.City,
		info.Gender,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserRepo) UpdateUserInfo(id int64, info UpdateUserInfo) error {
	_, err := r.conn.Exec("UPDATE user_info SET first_name=?, last_name=?, age=?, interests=?, city=?, gender=? WHERE user_id=?",
		info.FirstName,
		info.LastName,
		info.Age,
		info.Interests,
		info.City,
		info.Gender,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) AddFriends(id int64, friendIDs ...int64) error {
	// TODO: use transaction here

	for _, friendID := range friendIDs {
		_, err := r.conn.Exec("INSERT INTO friends (user_id, friend_id) VALUES (?, ?)", id, friendID)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *UserRepo) DeleteFriends(id int64, friendIDs ...int64) error {
	// TODO: use transaction here

	for _, friendID := range friendIDs {
		_, err := r.conn.Exec("DELETE FROM friends WHERE user_id=? AND friend_id=?", id, friendID)
		if err != nil {
			return err
		}
	}

	return nil
}
