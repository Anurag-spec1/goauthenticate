package controllers

import (
    "context"
    "time"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    
    "github.com/Anurag-spec1/goauthenticate/config"
    "github.com/Anurag-spec1/goauthenticate/models"
    "github.com/Anurag-spec1/goauthenticate/utils"
    "github.com/Anurag-spec1/goauthenticate/services"
)

func RequestOTP(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    // Validate college domain
    if !utils.ValidateCollegeDomain(req.Email) {
        c.JSON(400, gin.H{
            "success": false,
            "error": "Only @kiet.edu emails are allowed",
        })
        return
    }

    // Parse email information
    emailInfo := utils.ParseCollegeEmail(req.Email)
    if !emailInfo.IsValidFormat {
        c.JSON(400, gin.H{
            "success": false,
            "error": "Invalid college email format",
            "expected_format": "name.yyyybranchroll@kiet.edu",
            "example": "anurag.2428cse2059@kiet.edu",
        })
        return
    }

    // Generate OTP
    otp := utils.GenerateOTP()
    otpExpiresAt := time.Now().Add(10 * time.Minute)

    // Check if user exists
    var user models.User
    err := config.UserCollection.FindOne(
        context.Background(),
        bson.M{"email": req.Email},
    ).Decode(&user)

    if err != nil {
        // User doesn't exist, create new user
        if err == mongo.ErrNoDocuments {
            // Calculate current year based on admission year
            currentYear, yearNumber := calculateCurrentYearBasedOn2029(emailInfo.AdmissionYear)
            
            user = models.User{
                ID:            primitive.NewObjectID(),
                Name:          emailInfo.Name,
                Email:         req.Email,
                RollNumber:    emailInfo.RollNumber,
                Branch:        emailInfo.Branch,
                AdmissionYear: emailInfo.AdmissionYear,
                CurrentYear:   currentYear,
                YearNumber:    yearNumber,
                Batch:         emailInfo.Batch,
                OTP:           otp,
                OTPExpiresAt:  otpExpiresAt,
                IsVerified:    false,
                CreatedAt:     time.Now(),
            }
            
            _, err = config.UserCollection.InsertOne(context.Background(), user)
            if err != nil {
                c.JSON(500, gin.H{
                    "success": false,
                    "error": "Failed to create user",
                })
                return
            }
        } else {
            c.JSON(500, gin.H{
                "success": false,
                "error": "Database error",
            })
            return
        }
    } else {
        // Update existing user's OTP
        update := bson.M{
            "$set": bson.M{
                "otp":           otp,
                "otp_expires_at": otpExpiresAt,
            },
        }
        
        _, err = config.UserCollection.UpdateOne(
            context.Background(),
            bson.M{"email": req.Email},
            update,
        )
        if err != nil {
            c.JSON(500, gin.H{
                "success": false,
                "error": "Failed to update OTP",
            })
            return
        }
    }
	
    if err := services.SendOTPEmail(req.Email, otp); err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "error": "Failed to send OTP email",
        })
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "OTP sent successfully",
        "email": req.Email,
        "data_extracted": gin.H{
            "name":           emailInfo.Name,
            "roll_number":    emailInfo.RollNumber,
            "branch":         emailInfo.Branch,
            "admission_year": emailInfo.AdmissionYear,
            "current_year":   emailInfo.CurrentYear,
            "year_number":    emailInfo.YearNumber,
            "batch":          emailInfo.Batch,
        },
    })
}

// Helper function to calculate current year based on 2029 = 1st year
func calculateCurrentYearBasedOn2029(admissionYear string) (string, int) {
    yearInt, err := strconv.Atoi(admissionYear)
    if err != nil {
        return "1st Year", 1
    }
    
    // Fixed mapping based on your requirement:
    // 2029 = 1st year
    // 2028 = 2nd year  
    // 2027 = 3rd year
    // 2026 = 4th year
    
    // Simple formula: YearNumber = 2029 - AdmissionYear + 1
    yearNumber := 2029 - yearInt + 1
    
    // Ensure year is between 1 and 4
    if yearNumber < 1 {
        yearNumber = 1
    } else if yearNumber > 4 {
        yearNumber = 4
    }
    
    // Convert to string
    var yearString string
    switch yearNumber {
    case 1:
        yearString = "1st Year"
    case 2:
        yearString = "2nd Year"
    case 3:
        yearString = "3rd Year"
    case 4:
        yearString = "4th Year"
    default:
        yearString = "Graduated"
    }
    
    return yearString, yearNumber
}

func VerifyOTP(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required,email"`
        OTP   string `json:"otp" binding:"required,min=6,max=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "error": "Invalid request",
        })
        return
    }

    // Find user by email
    var user models.User
    err := config.UserCollection.FindOne(
        context.Background(),
        bson.M{"email": req.Email},
    ).Decode(&user)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(404, gin.H{
                "success": false,
                "error": "User not found",
            })
        } else {
            c.JSON(500, gin.H{
                "success": false,
                "error": "Database error",
            })
        }
        return
    }

    // Verify OTP
    if !utils.IsOTPValid(user.OTP, req.OTP, user.OTPExpiresAt) {
        c.JSON(401, gin.H{
            "success": false,
            "error": "Invalid or expired OTP",
        })
        return
    }

    // Clear OTP after successful verification
    update := bson.M{
        "$set": bson.M{
            "otp":         "",
            "is_verified": true,
        },
    }
    
    _, err = config.UserCollection.UpdateOne(
        context.Background(),
        bson.M{"email": req.Email},
        update,
    )
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "error": "Failed to update user",
        })
        return
    }

    // Generate JWT tokens
    accessToken, err := utils.GenerateAccessToken(user.ID.Hex())
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "error": "Failed to generate access token",
        })
        return
    }

    refreshToken, err := utils.GenerateRefreshToken(user.ID.Hex())
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "error": "Failed to generate refresh token",
        })
        return
    }

    // Store refresh token in database
    _, err = config.UserCollection.UpdateOne(
        context.Background(),
        bson.M{"email": req.Email},
        bson.M{"$set": bson.M{"refresh_token": refreshToken}},
    )
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "error": "Failed to store refresh token",
        })
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "Authentication successful",
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "user": gin.H{
            "id":             user.ID.Hex(),
            "name":           user.Name,
            "email":          user.Email,
            "roll_number":    user.RollNumber,
            "branch":         user.Branch,
            "admission_year": user.AdmissionYear,
            "current_year":   user.CurrentYear,
            "year_number":    user.YearNumber,
            "batch":          user.Batch,
        },
    })
}

func Refresh(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "error": "Invalid request",
        })
        return
    }

    // Parse and validate refresh token
    token, err := utils.ParseToken(req.RefreshToken, true)
    if err != nil || !token.Valid {
        c.JSON(401, gin.H{
            "success": false,
            "error": "Invalid refresh token",
        })
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(401, gin.H{
            "success": false,
            "error": "Invalid token claims",
        })
        return
    }

    userID, ok := claims["user_id"].(string)
    if !ok {
        c.JSON(401, gin.H{
            "success": false,
            "error": "Invalid user ID in token",
        })
        return
    }

    // Verify refresh token exists in database
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        c.JSON(401, gin.H{
            "success": false,
            "error": "Invalid user ID format",
        })
        return
    }

    var user models.User
    err = config.UserCollection.FindOne(
        context.Background(),
        bson.M{
            "_id":           objID,
            "refresh_token": req.RefreshToken,
        },
    ).Decode(&user)

    if err != nil {
        c.JSON(401, gin.H{
            "success": false,
            "error": "Refresh token not found or invalid",
        })
        return
    }

    // Generate new access token
    newAccessToken, err := utils.GenerateAccessToken(userID)
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "error": "Failed to generate access token",
        })
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "access_token": newAccessToken,
    })
}

func GetProfile(c *gin.Context) {
    // Get user ID from middleware
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, gin.H{
            "success": false,
            "error": "User not authenticated",
        })
        return
    }

    objID, err := primitive.ObjectIDFromHex(userID.(string))
    if err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "error": "Invalid user ID",
        })
        return
    }

    var user models.User
    err = config.UserCollection.FindOne(
        context.Background(),
        bson.M{"_id": objID},
    ).Decode(&user)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(404, gin.H{
                "success": false,
                "error": "User not found",
            })
        } else {
            c.JSON(500, gin.H{
                "success": false,
                "error": "Database error",
            })
        }
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "user": gin.H{
            "id":             user.ID.Hex(),
            "name":           user.Name,
            "email":          user.Email,
            "roll_number":    user.RollNumber,
            "branch":         user.Branch,
            "admission_year": user.AdmissionYear,
            "current_year":   user.CurrentYear,
            "year_number":    user.YearNumber,
            "batch":          user.Batch,
            "is_verified":    user.IsVerified,
            "created_at":     user.CreatedAt.Format(time.RFC3339),
        },
    })
}