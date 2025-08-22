package main

import (
	"fmt"

	// Import the local smtp package
	smtp "github.com/fajar/learn-go/04-smtp"
)

func main() {
	// Create a test configuration (no actual sending)
	config := smtp.EmailConfig{
		SMTPServer:   "localhost",
		SMTPPort:     25,
		SMTPUsername: "test",
		SMTPPassword: "test",
		SenderEmail:  "test@example.com",
		SenderName:   "Test Sender",
	}

	// Create a test message
	message := smtp.EmailMessage{
		To:        []string{"recipient@example.com"},
		Subject:   "Test Email",
		PlainBody: "This is a test email.",
	}

	// Print the configuration and message details to verify structure
	fmt.Println("=== SMTP Package Test ===")
	fmt.Println("This is a test of the SMTP package structure. No emails will be sent.")
	fmt.Println()

	fmt.Println("Email Configuration:")
	fmt.Printf("SMTP Server: %s:%d\n", config.SMTPServer, config.SMTPPort)
	fmt.Printf("Sender: %s <%s>\n", config.SenderName, config.SenderEmail)
	fmt.Println()

	fmt.Println("Email Message:")
	fmt.Printf("To: %v\n", message.To)
	fmt.Printf("Subject: %s\n", message.Subject)
	fmt.Printf("Body: %s\n", message.PlainBody)

	// Create an email sender (but don't actually send)
	_ = smtp.NewEmailSender(config)
	fmt.Println("\nSMTP package structure test completed successfully!")
}