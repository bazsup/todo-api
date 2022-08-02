package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/bazsup/todoapi/todo"
)

var nextId = 1

type LocalStore struct {
	tasks []todo.Todo
	mu    sync.Mutex
}

func NewLocalStore() *LocalStore {
	return &LocalStore{tasks: []todo.Todo{}, mu: sync.Mutex{}}
}

func (s *LocalStore) New(todo *todo.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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

func (s *LocalStore) Update(id int, completed bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id > len(s.tasks) || 0 > id {
		return fmt.Errorf("id not found")
	}

	s.tasks[id-1].Completed = completed
	return nil
}

func (s *LocalStore) Reset() {
	s.tasks = []todo.Todo{}
	nextId = 1
}
