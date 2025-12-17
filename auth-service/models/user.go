package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	Password     string             `json:"password"`
	RefreshToken string             `bson:"refresh_token,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
}