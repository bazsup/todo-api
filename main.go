package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bazsup/todoapi/auth"
	"github.com/bazsup/todoapi/todo"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("erfailed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/tokenz", auth.AccessToken)

	handler := todo.NewTodoHandler(db)
	r.POST("/todos", handler.NewTask)

	r.Run()
}
