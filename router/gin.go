package router

import (
	"github.com/bazsup/todoapi/todo"
	"github.com/gin-gonic/gin"
)

type MyContext struct {
	*gin.Context
}

func NewMyContext(c *gin.Context) *MyContext {
	return &MyContext{c}
}

func (c *MyContext) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}
func (c *MyContext) TransactionID() string {
	return c.Request.Header.Get("TransactionID")
}
func (c *MyContext) Audience() string {
	if aud, ok := c.Get("aud"); ok {
		if s, ok := aud.(string); ok {
			return s
		}
	}
	return ""
}
func (c *MyContext) JSON(statuscode int, v interface{}) {
	c.Context.JSON(statuscode, v)
}

func NewGinHandler(handler func(todo.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(NewMyContext(ctx))
	}
}