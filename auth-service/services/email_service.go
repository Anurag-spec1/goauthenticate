package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

type EmailService struct{}

func NewEmailService() *EmailService {
    return &EmailService{}
}

func (es *EmailService) SendOTPEmail(to, otp string) error {
    provider := os.Getenv("EMAIL_PROVIDER")
    
    // ADD THIS DEBUG LOGGING
    log.Printf("🔍 DEBUG: EMAIL_PROVIDER='%s'", provider)
    log.Printf("🔍 DEBUG: RESEND_API_KEY exists: %v", os.Getenv("RESEND_API_KEY") != "")
    log.Printf("🔍 DEBUG: EMAIL_FROM='%s'", os.Getenv("EMAIL_FROM"))
    
    if provider == "resend" {
        log.Println("🔍 DEBUG: Using Resend provider")
        return es.sendViaResend(to, otp)
    }
    
    log.Printf("🔍 DEBUG: Falling back to simulation (provider='%s')", provider)
    // Fallback to simulation
    return es.simulateEmail(to, otp)
}

func (es *EmailService) sendViaResend(to, otp string) error {
    apiKey := os.Getenv("RESEND_API_KEY")
    from := os.Getenv("EMAIL_FROM")
    
    if apiKey == "" {
        log.Println("⚠️ RESEND_API_KEY not set, falling back to simulation")
        return es.simulateEmail(to, otp)
    }
    
    if from == "" {
        from = "onboarding@resend.dev"
    }
    
    // HTML email template
    htmlContent := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>KIET Authentication OTP</title>
        <style>
            body { font-family: 'Arial', sans-serif; background-color: #f7f9fc; margin: 0; padding: 0; }
            .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 20px rgba(0,0,0,0.1); }
            .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center; color: white; }
            .content { padding: 40px; }
            .otp-container { background: #f8f9fa; border: 2px dashed #667eea; border-radius: 10px; padding: 25px; text-align: center; margin: 30px 0; }
            .otp-code { font-size: 42px; font-weight: bold; color: #667eea; letter-spacing: 10px; font-family: 'Courier New', monospace; margin: 15px 0; }
            .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; text-align: center; }
            .security-note { background: #e3f2fd; padding: 15px; border-radius: 8px; margin: 20px 0; border-left: 4px solid #2196f3; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1 style="margin: 0; font-size: 28px;">🔐 KIET Authentication</h1>
                <p style="margin: 10px 0 0; opacity: 0.9;">Secure Access Portal</p>
            </div>
            <div class="content">
                <h2 style="color: #333; margin-bottom: 20px;">Hello KIET Student!</h2>
                <p style="color: #555; line-height: 1.6;">Your One-Time Password (OTP) for authentication is:</p>
                
                <div class="otp-container">
                    <div class="otp-code">%s</div>
                    <p style="color: #666; margin: 10px 0; font-size: 14px;">⏱️ Valid for 10 minutes</p>
                </div>
                
                <div class="security-note">
                    <strong style="color: #1976d2;">🔒 Security Notice:</strong>
                    <ul style="margin: 10px 0; padding-left: 20px; color: #555;">
                        <li>Do NOT share this OTP with anyone</li>
                        <li>KIET staff will never ask for your OTP</li>
                        <li>If you didn't request this, please ignore this email</li>
                    </ul>
                </div>
                
                <p style="color: #555; line-height: 1.6;">Enter this OTP in the authentication portal to complete your login.</p>
                
                <div class="footer">
                    <p style="margin: 5px 0;"><strong>KIET Group of Institutions</strong></p>
                    <p style="margin: 5px 0; color: #888;">Delhi-NCR, Ghaziabad, Uttar Pradesh</p>
                    <p style="margin: 10px 0; font-size: 11px; color: #999;">This is an automated message. Please do not reply.</p>
                    <p style="margin: 5px 0; font-size: 11px; color: #999;">Time: %s</p>
                </div>
            </div>
        </div>
    </body>
    </html>
    `, otp, time.Now().Format("2006-01-02 15:04:05"))
    
    // Plain text version for email clients that don't support HTML
    textContent := fmt.Sprintf(`
KIET Authentication System
==========================

Your One-Time Password (OTP) is: %s

This OTP is valid for 10 minutes.

SECURITY NOTICE:
• Do NOT share this OTP with anyone
• KIET staff will never ask for your OTP
• If you didn't request this, please ignore this email

Enter this OTP in the authentication portal to complete your login.

---
KIET Group of Institutions
Delhi-NCR, Ghaziabad, Uttar Pradesh

This is an automated message. Please do not reply.
Time: %s
`, otp, time.Now().Format("2006-01-02 15:04:05"))
    
    // Create Resend API request
    payload := map[string]interface{}{
        "from":    "KIET Authentication <" + from + ">",
        "to":      []string{to},
        "subject": "Your KIET Authentication OTP",
        "html":    htmlContent,
        "text":    textContent,
        "reply_to": "no-reply@kiet.edu",
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        log.Printf("❌ Error marshaling email data: %v", err)
        return es.simulateEmail(to, otp)
    }
    
    // Send request to Resend API
    req, err := http.NewRequest("POST", "https://api.resend.com/emails", 
                                bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("❌ Error creating HTTP request: %v", err)
        return es.simulateEmail(to, otp)
    }
    
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("❌ Error sending email via Resend: %v", err)
        return es.simulateEmail(to, otp)
    }
    defer resp.Body.Close()
    
    // Parse response
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Printf("❌ Error parsing Resend response: %v", err)
    }
    
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        log.Printf("✅ Email sent successfully via Resend to: %s", to)
        log.Printf("   Resend ID: %v", result["id"])
        return nil
    }
    
    // Resend API error
    log.Printf("❌ Resend API error (Status: %d): %v", resp.StatusCode, result)
    return es.simulateEmail(to, otp)
}

func (es *EmailService) simulateEmail(to, otp string) error {
    log.Printf("📧 [SIMULATION] OTP for %s: %s", to, otp)
    
    border := strings.Repeat("═", 60)
    
    fmt.Printf("\n%s\n", border)
    fmt.Println("📧 EMAIL SIMULATION MODE")
    fmt.Println(border)
    fmt.Printf("To: %s\n", to)
    fmt.Printf("OTP: %s\n", otp)
    fmt.Printf("Valid until: %s\n", time.Now().Add(10*time.Minute).Format("15:04:05"))
    fmt.Println(border + "\n")
    
    return nil
}