package main

import (
	"fmt"
	"os"
	"strconv"

	authdl "otus-hl-network/internal/auth/delivery"
	"otus-hl-network/internal/domain"
	"otus-hl-network/internal/server"
	userdl "otus-hl-network/internal/user/delivery"
	userrepo "otus-hl-network/internal/user/repository"
	userua "otus-hl-network/internal/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type NilAuthManager struct{}

func (n NilAuthManager) Generate(user domain.User) (domain.UserContext, error) {
	return domain.UserContext{
		ID:       1,
		Username: "foo@gmail.com",
		Token:    "",
	}, nil
}

func (n NilAuthManager) Verify(accessToken string) (domain.UserContext, error) {
	return domain.UserContext{
		ID:       1,
		Username: "foo@gmail.com",
		Token:    "",
	}, nil
}

func main() {
	sqlConn, err := sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOSTNAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE")))
	if err != nil {
		logrus.WithError(err).Error("could not connect to DB")
		return
	}

	var backendPort int
	portStr := os.Getenv("PORT")
	if x, err := strconv.Atoi(portStr); err == nil {
		backendPort = x
	} else {
		backendPort = 8080
	}

	// jwtManager := auth.NewJWTManager(os.Getenv("JWT_SECRET_KEY"), time.Hour*24)
	jwtManager := &NilAuthManager{}

	// user initialization
	userRepo := userrepo.NewUserRepo(sqlConn)
	userUA := userua.NewUsecase(userRepo)
	userDelievery := userdl.New(userUA)

	// auth initialization
	authDelivery := authdl.New(jwtManager, userUA)

	srv := server.New(jwtManager,
		userDelievery,
		authDelivery,
	)

	if err := srv.Run(backendPort); err != nil {
		logrus.WithError(err).Fatal("web-server stopped with error")
	}
}
