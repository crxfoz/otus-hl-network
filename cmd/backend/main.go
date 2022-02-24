package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/crxfoz/otus-hl-network/internal/auth"
	authdl "github.com/crxfoz/otus-hl-network/internal/auth/delivery"
	chatdl "github.com/crxfoz/otus-hl-network/internal/chat/delivery"
	chatrepo "github.com/crxfoz/otus-hl-network/internal/chat/repository"
	"github.com/crxfoz/otus-hl-network/internal/domain"
	"github.com/crxfoz/otus-hl-network/internal/server"
	userdl "github.com/crxfoz/otus-hl-network/internal/user/delivery"
	userrepo "github.com/crxfoz/otus-hl-network/internal/user/repository"
	userua "github.com/crxfoz/otus-hl-network/internal/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mongoConn, err := mongo.NewClient(options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s/",
			os.Getenv("MONGO_HOST"),
			os.Getenv("MONGO_PORT"),
		),
	))
	if err != nil {
		logrus.WithError(err).Error("could not create mongo")
		return
	}

	if err := mongoConn.Connect(context.Background()); err != nil {
		logrus.WithError(err).Error("could not connect to mongo")
		return
	}

	mDB := mongoConn.Database("chat")

	var jwtManager auth.AuthManager

	if os.Getenv("DEBUG") == "1" {
		// use it for testing
		// it disables auth and it may cause some problems with testing cuz it says that user is auth'ed as user with id=1
		jwtManager = &NilAuthManager{}
	} else {
		// use it for _prod_
		jwtManager = auth.NewJWTManager(os.Getenv("JWT_SECRET_KEY"), time.Hour*24)
	}

	// user initialization
	userRepo := userrepo.NewUserRepo(sqlConn)
	userUA := userua.NewUsecase(userRepo)
	userDelievery := userdl.New(userUA)

	// auth initialization
	authDelivery := authdl.New(jwtManager, userUA)

	// chat initialization
	chatRepo := chatrepo.NewMongo(context.Background(), mDB, time.Second*5)
	chatDelievery := chatdl.NewChatHandler(chatRepo)

	srv := server.New(jwtManager,
		userDelievery,
		authDelivery,
		chatDelievery,
	)

	if err := srv.Run(backendPort); err != nil {
		logrus.WithError(err).Fatal("web-server stopped with error")
	}
}
