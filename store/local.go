package store

import (
	"time"

	"github.com/bazsup/todoapi/todo"
)

var nextId = 1

type LocalStore struct {
	tasks []todo.Todo
}

func NewLocalStore() *LocalStore {
	return &LocalStore{tasks: []todo.Todo{}}
}

func (s *LocalStore) New(todo *todo.Todo) error {
	todo.ID = uint(nextId)
	now := time.Now().UTC()
	todo.CreatedAt = now
	todo.UpdatedAt = now
	nextId++
	s.tasks = append(s.tasks, *todo)
	return nil
}

func (s *LocalStore) GetAll() ([]todo.Todo, error) {
	return s.tasks, nil
}
