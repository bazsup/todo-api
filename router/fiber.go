package router

import (
	"github.com/bazsup/todoapi/todo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{c}
}

func (c *FiberCtx) Bind(v interface{}) error {
	return c.Ctx.BodyParser(v)
}
func (c *FiberCtx) TransactionID() string {
	return string(c.Request().Header.Peek("TransactionID"))
}
func (c *FiberCtx) Audience() string {
	return c.Ctx.Get("aud")
}
func (c *FiberCtx) JSON(statuscode int, v interface{}) {
	c.Ctx.Status(statuscode).JSON(v)
}
func (c *FiberCtx) Authorization() string {
	auth := string(c.Request().Header.Peek("Authorization"))
	return auth
}

func (c *FiberCtx) AbortWithStatus(statuscode int) {
	c.Status(statuscode)
}
func (c *FiberCtx) Set(k string, v interface{}) {
	if s, ok := v.(string); ok {
		c.Ctx.Set(k, s)
	}
}
func (c *FiberCtx) Next() {
	c.Ctx.Next()
}

func NewFiberHandler(handler func(todo.Context)) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		handler(NewFiberCtx(ctx))
		return nil
	}
}

type FiberRouter struct {
	*fiber.App
}

func NewFiberRouter() *FiberRouter {
	r := fiber.New()

	r.Use(cors.New())
	r.Use(logger.New())

	return &FiberRouter{r}
}

func (r *FiberRouter) GET(path string, handler func(todo.Context)) {
	r.App.Get(path, NewFiberHandler(handler))
}

func (r *FiberRouter) POST(path string, handler func(todo.Context)) {
	r.App.Post(path, func(c *fiber.Ctx) error {
		handler(NewFiberCtx(c))
		return nil
	})
}

type FiberRouterGroup struct {
	g *fiber.Group
}

func (r *FiberRouter) Group(path string, handler func(todo.Context)) *FiberRouterGroup {
	return &FiberRouterGroup{r.App.Group(path, NewFiberHandler(handler)).(*fiber.Group)}
}

func (r *FiberRouterGroup) POST(path string, handler func(todo.Context)) {
	r.g.Post(path, NewFiberHandler(handler))
}
