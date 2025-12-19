package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/Anurag-spec1/goauthenticate/utils"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{
                "success": false,
                "error": "Authorization header is required",
            })
            c.Abort()
            return
        }

        // Remove "Bearer " prefix
        token = strings.TrimPrefix(token, "Bearer ")
        token = strings.TrimSpace(token)

        if token == "" {
            c.JSON(401, gin.H{
                "success": false,
                "error": "Token is empty",
            })
            c.Abort()
            return
        }

        // Parse and validate token
        parsedToken, err := utils.ParseToken(token, false)
        if err != nil || !parsedToken.Valid {
            c.JSON(401, gin.H{
                "success": false,
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        // Extract claims
        claims, ok := parsedToken.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(401, gin.H{
                "success": false,
                "error": "Invalid token claims",
            })
            c.Abort()
            return
        }

        // Get user ID from claims
        userID, ok := claims["user_id"].(string)
        if !ok || userID == "" {
            c.JSON(401, gin.H{
                "success": false,
                "error": "Invalid user ID in token",
            })
            c.Abort()
            return
        }

        // Set user ID in context for use in controllers
        c.Set("user_id", userID)
        c.Next()
    }
}