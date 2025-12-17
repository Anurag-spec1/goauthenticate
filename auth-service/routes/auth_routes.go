package routes

import (
	"github.com/Anurag-spec1/goauthenticate/controllers"
	"github.com/Anurag-spec1/goauthenticate/middleware"

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
