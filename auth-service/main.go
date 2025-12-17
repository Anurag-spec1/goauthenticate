package main

import (
	"auth-service/config"
	"auth-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	r := gin.Default()
	routes.RegisterAuthRoutes(r)

	r.Run(":8080")
}
