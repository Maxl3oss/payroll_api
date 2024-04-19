package utils

import (
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
)

func GenerateOTP() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

// Function to send OTP to email using Gmail SMTP
func SendEmail(email, otp string) error {
	// Replace these values with your Gmail credentials
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUsername := os.Getenv("SEND_MAIL_USER")
	smtpPassword := os.Getenv("SEND_MAIL_PASS")

	// Compose email message
	message := "From: " + smtpUsername + "\n" +
		"To: " + email + "\n" +
		"Subject: Password Reset OTP\n\n" +
		"Your OTP for password reset is: " + otp

	// Send email using Gmail SMTP
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUsername, []string{email}, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
