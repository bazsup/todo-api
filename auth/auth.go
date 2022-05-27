package auth

import (
	"net/http"
	"time"

	"github.com/bazsup/todoapi/todo"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AccessToken(signature []byte) func(c todo.Context) {
	return func(c todo.Context) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(),
			Audience:  "Bas",
		})

		ss, err := token.SignedString(signature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": ss})
	}
}
