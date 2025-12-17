package main

import (
	"github.com/Anurag-spec1/goauthenticate/config"
	"github.com/Anurag-spec1/goauthenticate/routes"

	"github.com/gin-gonic/gin"
)


func main() {
	config.LoadEnv()
	config.ConnectDB()

	r := gin.Default()
	routes.RegisterAuthRoutes(r)

	r.Run(":8080")
}
