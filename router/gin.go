package router

import (
	"github.com/bazsup/todoapi/todo"
	"github.com/gin-contrib/cors"
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

func (c *MyContext) Authorization() string {
	auth := c.Request.Header.Get("Authorization")
	return auth
}
func (c *MyContext) AbortWithStatus(statuscode int) {
	c.Context.AbortWithStatus(statuscode)
}
func (c *MyContext) Set(k string, v interface{}) {
	c.Context.Set(k, v)
}
func (c *MyContext) Next() {
	c.Context.Next()
}

func (c *MyContext) Param(key string) string {
	return c.Context.Param(key)
}

func NewGinHandler(handler func(todo.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(NewMyContext(ctx))
	}
}

type MyRouter struct {
	*gin.Engine
}

func NewMyRouter() *MyRouter {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}
	r.Use(cors.New(config))

	return &MyRouter{r}
}

func (r *MyRouter) POST(path string, handler func(todo.Context)) {
	r.Engine.POST(path, NewGinHandler(handler))
}

type MyRouterGroup struct {
	*gin.RouterGroup
}

func (r *MyRouter) Group(path string, handler func(todo.Context)) *MyRouterGroup {	
	return &MyRouterGroup{r.Engine.Group(path, NewGinHandler(handler))}
}

func (r *MyRouterGroup) POST(path string, handler func(todo.Context)) {
	r.RouterGroup.POST(path, NewGinHandler(handler))
}
