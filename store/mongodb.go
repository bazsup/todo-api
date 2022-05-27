package store

import (
	"context"

	"github.com/bazsup/todoapi/todo"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStore struct {
	*mongo.Collection
}

func NewMongoDBStore(collection *mongo.Collection) *MongoDBStore {
	return &MongoDBStore{collection}
}

func (s *MongoDBStore) New(todo *todo.Todo) error {
	_, err := s.Collection.InsertOne(context.Background(), todo)
	return err
}
