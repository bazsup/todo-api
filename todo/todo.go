package todo

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type Todo struct {
	Title     string `json:"text" binding:"required"`
	ID        uint   `gorm:"primary_key"`
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Todo) TableName() string {
	return "todos"
}

type storer interface {
	New(*Todo) error
	GetAll() ([]Todo, error)
	Update(int, bool) error
	Reset()
}

type TodoHandler struct {
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store}
}

type Context interface {
	Param(string) string
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

type UpdateStatusReq struct {
	Completed bool `json:"completed" binding:"required"`
}

func (t *TodoHandler) UpdateStatus(c Context) {

	var req UpdateStatusReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	err = t.store.Update(id, req.Completed)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"data": "ok"})
}

func (t *TodoHandler) Reset(c Context) {
	t.store.Reset()

	c.JSON(http.StatusOK, map[string]interface{}{"data": "ok"})
}

func (t *TodoHandler) GetTasks(c Context) {
	todos, err := t.store.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": todos,
	})
}
