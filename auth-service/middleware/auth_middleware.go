package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(401, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		token = strings.Replace(token, "Bearer ", "", 1)

		parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})

		if err != nil || !parsed.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
