package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Protect(signature []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				// return nil, jwt.ErrSignatureInvalid
			}

			return signature, nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			aud := claims["aud"]
			c.Set("aud", aud)
		}

		c.Next()
	}
}
