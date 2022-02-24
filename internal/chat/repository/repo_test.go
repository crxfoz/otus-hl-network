package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/crxfoz/otus-hl-network/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: write integration tests

var (
	once sync.Once
	repo *MongoSvc
)

func setup(t *testing.T) *MongoSvc {
	once.Do(func() {
		mongoConn, err := mongo.NewClient(options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s:%s/",
				"0.0.0.0",
				"27017",
			),
		))
		assert.Nil(t, err)

		err = mongoConn.Connect(context.Background())
		assert.Nil(t, err)

		mDB := mongoConn.Database("chat")

		repo = NewMongo(context.Background(), mDB, time.Second*10)
	})

	return repo
}

func TestMongoSvc_CreateChat(t *testing.T) {
	conn := setup(t)

	err := conn.CreateChat(uuid.New().String(), "", []int64{1, 2, 3})
	assert.Nil(t, err)
}

func TestMongoSvc_GetChat(t *testing.T) {
	conn := setup(t)

	chat, err := conn.GetChat("73fc670d-c6e7-4e42-aa46-3d16aaa70032b63")
	assert.Nil(t, err)

	fmt.Println(chat)
}

func TestMongoSvc_SendMessage(t *testing.T) {
	conn := setup(t)

	err := conn.SendMessage("73fc670d-c6e7-4e42-aa46-3d1670032b63", domain.Message{
		UserName: "kkk",
		Body:     "pewpew",
		// CreatedAt: time.Now().UTC().UnixMilli(),
		CreatedAt: 3,
	})
	assert.Nil(t, err)
}

func TestMongoSvc_GetHistory(t *testing.T) {
	conn := setup(t)

	msgs, err := conn.GetHistory("73fc670d-c6e7-4e42-aa46-3d1670032b63", 1000)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", msgs)
}
