package services

import (
    "fmt"
    "log"
    "net/smtp"
    "os"
)

func SendOTPEmail(to, otp string) error {
    from := os.Getenv("GMAIL_ADDRESS")
    password := os.Getenv("GMAIL_APP_PASSWORD")
    
    // If no credentials, just show OTP in console
    if from == "" || password == "" {
        log.Printf("📧 [SIMULATION] OTP for %s: %s", to, otp)
        fmt.Printf("\n=== OTP: %s (for %s) ===\n", otp, to)
        return nil
    }
    
    // Send real email
    auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
    
    msg := fmt.Sprintf("From: KIET Auth <%s>\r\n", from) +
           fmt.Sprintf("To: %s\r\n", to) +
           "Subject: Your KIET Authentication OTP\r\n" +
           "\r\n" +
           fmt.Sprintf("Your OTP is: %s\r\n", otp) +
           "Valid for 10 minutes.\r\n" +
           "Do not share with anyone.\r\n"
    
    err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
    
    if err != nil {
        log.Printf("❌ Email failed: %v", err)
        log.Printf("📧 [FALLBACK] OTP for %s: %s", to, otp)
        fmt.Printf("\n=== OTP (Email failed): %s ===\n", otp)
        return nil
    }
    
    log.Printf("✅ Email sent to: %s", to)
    return nil
}