package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name          string             `json:"name" bson:"name"`
    Email         string             `json:"email" bson:"email"`
    RollNumber    string             `json:"roll_number" bson:"roll_number"`
    Branch        string             `json:"branch" bson:"branch"`
    AdmissionYear string             `json:"admission_year" bson:"admission_year"`
    CurrentYear   string             `json:"current_year" bson:"current_year"`
    YearNumber    int                `json:"year_number" bson:"year_number"`
    Batch         string             `json:"batch" bson:"batch"`
    OTP           string             `json:"-" bson:"otp,omitempty"`
    OTPExpiresAt  time.Time          `json:"-" bson:"otp_expires_at,omitempty"`
    RefreshToken  string             `json:"-" bson:"refresh_token,omitempty"`
    IsVerified    bool               `json:"is_verified" bson:"is_verified"`
    CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}