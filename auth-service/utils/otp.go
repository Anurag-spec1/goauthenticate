package utils

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "time"
)

func GenerateOTP() string {
    // Generate 6-digit OTP
    max := big.NewInt(1000000)
    n, err := rand.Int(rand.Reader, max)
    if err != nil {
        // Fallback to timestamp-based OTP
        return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
    }
    return fmt.Sprintf("%06d", n.Int64())
}

func IsOTPValid(storedOTP, providedOTP string, expiresAt time.Time) bool {
    if storedOTP == "" || providedOTP == "" {
        return false
    }
    if storedOTP != providedOTP {
        return false
    }
    return time.Now().Before(expiresAt)
}