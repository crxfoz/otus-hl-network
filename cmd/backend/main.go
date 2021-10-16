package main

import (
	"fmt"
	"os"
	"time"

	"otus-hl-network/internal/auth"
	"otus-hl-network/internal/server"
	userdl "otus-hl-network/internal/user/delivery"
	userrepo "otus-hl-network/internal/user/repository"
	userua "otus-hl-network/internal/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func main() {

	sqlConn, err := sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@(%s:3306)/%s",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOSTNAME"),
			os.Getenv("DB_DATABASE")))
	if err != nil {
		logrus.WithError(err).Error("could not connect to DB")
		return
	}

	jwtManager := auth.NewJWTManager(os.Getenv("JWT_SECRET_KEY"), time.Hour*24)

	// user initialization
	userRepo := userrepo.NewUserRepo(sqlConn)
	userUA := userua.NewUsecase(userRepo)
	userDelievery := userdl.New(userUA)

	srv := server.New(jwtManager,
		userDelievery,
	)

	if err := srv.Run(8080); err != nil {
		logrus.WithError(err).Fatal("web-server stopped with error")
	}
}
