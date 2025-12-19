package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anurag-spec1/goauthenticate/config"
	"github.com/Anurag-spec1/goauthenticate/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Set Gin mode
    if os.Getenv("GIN_MODE") == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Connect to MongoDB
    config.ConnectDB()
    defer config.DisconnectDB()

    // Setup Gin router with middleware
    r := gin.Default()
    
    // Add CORS middleware for development
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // Register routes
    routes.RegisterAuthRoutes(r)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Graceful shutdown
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }

    // Start server in a goroutine
    go func() {
        fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
        fmt.Println("📧 OTP Authentication API for KIET College")
        fmt.Println("==========================================")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    // Give server 5 seconds to finish current requests
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exited properly")
}