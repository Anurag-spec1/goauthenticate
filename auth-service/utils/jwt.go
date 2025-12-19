package utils

import (
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID string) (string, error) {
    secret := os.Getenv("ACCESS_SECRET")
    if secret == "" {
        secret = "myaccesssecret"
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(15 * time.Minute).Unix(),
        "iat":     time.Now().Unix(),
        "type":    "access",
    })

    return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID string) (string, error) {
    secret := os.Getenv("REFRESH_SECRET")
    if secret == "" {
        secret = "myrefreshsecret"
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
        "iat":     time.Now().Unix(),
        "type":    "refresh",
    })

    return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, isRefresh bool) (*jwt.Token, error) {
    var secret string
    if isRefresh {
        secret = os.Getenv("REFRESH_SECRET")
        if secret == "" {
            secret = "myrefreshsecret"
        }
    } else {
        secret = os.Getenv("ACCESS_SECRET")
        if secret == "" {
            secret = "myaccesssecret"
        }
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(secret), nil
    })

    return token, err
}

func ExtractUserIDFromToken(tokenString string, isRefresh bool) (string, error) {
    token, err := ParseToken(tokenString, isRefresh)
    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if userID, ok := claims["user_id"].(string); ok {
            return userID, nil
        }
    }
    
    return "", jwt.ErrInvalidKey
}