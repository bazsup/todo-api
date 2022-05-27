package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/bazsup/todoapi/auth"
	"github.com/bazsup/todoapi/router"
	"github.com/bazsup/todoapi/store"
	"github.com/bazsup/todoapi/todo"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Println("please consider enviroment variables: %s", err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	// r := router.NewMyRouter()
	r := router.NewFiberRouter()

	// r.GET("/healthz", func(c *gin.Context) {
	// 	c.Status(200)
	// })
	// r.GET("/limitz", limitedHandler)
	// r.GET("/x", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"buildcommit": buildcommit,
	// 		"buildtime":   buildtime,
	// 	})
	// })
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	r.GET("/tokenz", auth.AccessToken([]byte(os.Getenv("SIGN"))))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_CONN")))
	if err != nil {
		panic("failed to connect mongo")
	}
	collection := client.Database("myapp").Collection("todos")

	mongoStore := store.NewMongoDBStore(collection)
	// gormStore := store.NewGormStore(db)
	handler := todo.NewTodoHandler(mongoStore)

	protected.POST("/todos", handler.NewTask)

	if err := r.Listen(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// s := &http.Server{
	// 	Addr:           ":" + os.Getenv("PORT"),
	// 	Handler:        r,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// go func() {
	// 	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()

	// <-ctx.Done()
	// stop()
	// fmt.Println("shutting down gracefully, press ctrl+c again to force")

	// timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := s.Shutdown(timeoutCtx); err != nil {
	// 	fmt.Println(err)
	// }
}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
