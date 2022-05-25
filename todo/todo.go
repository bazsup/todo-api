package todo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bazsup/todoapi/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text" binding:"required"`
	gorm.Model
}

func (Todo) TableName() string {
	return "todos"
}

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db}
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	s := c.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(s, "Bearer ")

	fmt.Println(tokenString)
	if err := auth.Protect(tokenString); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}
