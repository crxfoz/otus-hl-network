package repository

import "gopkg.in/guregu/null.v4"

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserInfo struct {
	ID        int64       `db:"id"`
	UserID    int64       `db:"user_id"`
	FirstName null.String `db:"first_name"`
	LastName  null.String `db:"last_name"`
	Age       null.Int    `db:"age"`
	City      null.String `db:"city"`
	Interests null.String `db:"interests"` // TODO: replace with array?
	Gender    null.String `db:"gender"`
}

type UpdateUserInfo struct {
	FirstName string
	LastName  string
	Age       int
	City      string
	Interests string
	Gender    string
}
