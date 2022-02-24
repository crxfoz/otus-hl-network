package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/crxfoz/otus-hl-network/internal/domain"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSvc struct {
	conn      *mongo.Database
	ctx       context.Context
	opTimeout time.Duration
}

func NewMongo(ctx context.Context, conn *mongo.Database, opTimeout time.Duration) *MongoSvc {
	return &MongoSvc{
		ctx:       ctx,
		conn:      conn,
		opTimeout: opTimeout,
	}
}

func (m *MongoSvc) chats() *mongo.Collection {
	return m.conn.Collection("chats")
}

func (m *MongoSvc) messages() *mongo.Collection {
	return m.conn.Collection("messages")
}

func (m *MongoSvc) newContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(m.ctx, timeout)
}

func (m *MongoSvc) GetChat(id string) (domain.Chat, error) {
	res := m.chats().FindOne(m.ctx, bson.D{primitive.E{
		Key:   "id",
		Value: id,
	}})

	err := res.Err()

	if res.Err() != nil {
		return domain.Chat{}, fmt.Errorf("could not find chat: %v", err)
	}

	var chat domain.Chat

	if err := res.Decode(&chat); err != nil {
		return domain.Chat{}, fmt.Errorf("could not decode chat: %v", err)
	}

	return chat, nil
}

func (m *MongoSvc) CreateChat(id string, name string, users []int64) error {
	_, err := m.chats().InsertOne(m.ctx, domain.Chat{
		ID:    id,
		Users: users,
	})

	if err != nil {
		return fmt.Errorf("could not insert new chat: %v", err)
	}

	return nil
}

// GetHistory returns list of last messages
func (m *MongoSvc) GetHistory(id string, limit int) (domain.Messages, error) {
	fOpts := &options.FindOptions{}
	fOpts.SetSort(bson.D{{Key: "msg.createdat", Value: -1}})
	fOpts.SetLimit(int64(limit))

	cur, err := m.messages().Find(m.ctx, bson.D{primitive.E{Key: "id", Value: id}}, fOpts)
	if err != nil {
		return nil, fmt.Errorf("could not find items: %v", err)
	}

	out := make([]domain.Message, 0, limit)

	for cur.Next(m.ctx) {
		var currMsg chatMessage

		if err := cur.Decode(&currMsg); err != nil {
			logrus.WithError(err).Error("could not decode message")
			continue
		}

		out = append(out, currMsg.Msg)
	}

	return out, nil
}

type chatMessage struct {
	ID  string         `json:"id"`
	Msg domain.Message `json:"msg"`
}

func (m *MongoSvc) SendMessage(id string, msg domain.Message) error {
	_, err := m.messages().InsertOne(m.ctx, chatMessage{
		ID:  id,
		Msg: msg,
	})

	if err != nil {
		return fmt.Errorf("could not insert msg: %v", err)
	}

	return nil
}
