package routes

import (
	"auth-service/controllers"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.Refresh)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "This is protected profile"})
		})
	}
}
