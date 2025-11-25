package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	Jwt "sample_project/internal/pkg/jwt"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		claims, err := Jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
