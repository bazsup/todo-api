package todo_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bazsup/todoapi/todo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestNewTodoNotAllowSleep(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	payload := bytes.NewBufferString(`{"text":"sleep"}`)
	req, _ := http.NewRequest("POST", "/todos", payload)
	c.Request = req

	gormStore := todo.NewGormStore(&gorm.DB{})
	handler := todo.NewTodoHandler(gormStore)
	handler.NewTask(c)

	fmt.Println(w.Body.String())
}
