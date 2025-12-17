package controllers

import (
	"context"
	"os"
	"time"

	"auth-service/config"
	"auth-service/models"
	"auth-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)


func Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&req)

	hash, _ := utils.HashPassword(req.Password)

	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hash,
		CreatedAt: time.Now(),
	}

	_, err := config.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(201, gin.H{"message": "Registered successfully"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&req)

	var user models.User
	err := config.UserCollection.FindOne(
		context.Background(),
		bson.M{"email": req.Email},
	).Decode(&user)

	if err != nil || utils.CheckPassword(user.Password, req.Password) != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	access, _ := utils.GenerateAccessToken(user.ID.Hex())
	refresh, _ := utils.GenerateRefreshToken(user.ID.Hex())

	config.UserCollection.UpdateByID(
		context.Background(),
		user.ID,
		bson.M{"$set": bson.M{"refresh_token": refresh}},
	)

	c.JSON(200, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	c.ShouldBindJSON(&req)

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	newAccess, _ := utils.GenerateAccessToken(userID)

	c.JSON(200, gin.H{"access_token": newAccess})
}