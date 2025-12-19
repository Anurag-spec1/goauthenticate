package routes

import (
    "github.com/Anurag-spec1/goauthenticate/controllers"
    "github.com/Anurag-spec1/goauthenticate/middleware"
    "github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
    // Public routes
    r.POST("/auth/request-otp", controllers.RequestOTP)
    r.POST("/auth/verify-otp", controllers.VerifyOTP)
    r.POST("/auth/refresh", controllers.Refresh)

    // Protected routes (require authentication)
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", controllers.GetProfile)
        protected.GET("/test", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "This is a protected route",
                "success": true,
            })
        })
    }

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "OK",
            "message": "Authentication API is running",
        })
    })
}