package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func main() {
	sqlConn, err := sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s",
			"root",
			"user123",
			"0.0.0.0",
			"4000",
			"hl_network"))
	if err != nil {
		logrus.WithError(err).Error("could not connect to DB")
		return
	}

	successCnt := 0

	defer func() {
		logrus.WithField("num", successCnt).Info("wrote items in total")
	}()

	for i := 0; i < 1000000; i++ {
		name := fmt.Sprintf("foo-%d", i)
		age := i

		logrus.WithField("name", name).WithField("age", age).Info("trying to write")

		_, err := sqlConn.Exec("INSERT INTO example(name, age) VALUES(?, ?)", name, age)
		if err != nil {
			logrus.WithError(err).Error("could not write to DB")
			return
		}

		successCnt++
		logrus.WithField("name", name).WithField("age", age).Info("wrote successfully")
	}
}
