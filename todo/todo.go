package todo

import (
	"log"
	"net/http"
	"time"
)

type Todo struct {
	Title     string `json:"text" binding:"required"`
	ID        uint   `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Todo) TableName() string {
	return "todos"
}

type storer interface {
	New(*Todo) error
}

type TodoHandler struct {
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store}
}

type Context interface {
	Bind(interface{}) error
	TransactionID() string
	Audience() string
	JSON(int, interface{})
	Authorization() string
	AbortWithStatus(statuscode int)
	Set(k string, v interface{})
	Next()
}

func (t *TodoHandler) NewTask(c Context) {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	if todo.Title == "sleep" {
		transactionID := c.TransactionID()
		aud := c.Audience()
		log.Println(transactionID, aud, "not allowed")
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "not allowed"})
		return
	}

	err := t.store.New(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"ID": todo.ID,
	})
}
