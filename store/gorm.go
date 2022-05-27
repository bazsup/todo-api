package store

import (
	"github.com/bazsup/todoapi/todo"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db}
}

func (s *GormStore) New(todo *todo.Todo) error {
	return s.db.Create(todo).Error
}
