package main

import (
	"fmt"

	"github.com/fajar/learn-go/04-smtp"
)

func main() {
	fmt.Println("=== SMTP Email Sender Demo ===")
	fmt.Println("This is a demonstration of the SMTP email package.")
	fmt.Println("No emails will be sent in this demo.")
	fmt.Println()

	// Create email configuration with provided details
	config := smtp.EmailConfig{
		SMTPServer:         "mx.mailspace.id",
		SMTPPort:           587, // Using STARTTLS instead of SMTPS
		SMTPUsername:       "no-reply@ganapatih.com",
		SMTPPassword:       "risqor-0Wojsy-xipxenusetls", // Updated password
		SenderEmail:        "no-reply@ganapatih.com",
		SenderName:         "Go SMTP Test",
		InsecureSkipVerify: true, // Skip TLS verification for testing
		DebugMode:          true, // Enable debug logging
		AuthMethod:         "login", // Try LOGIN authentication method
	}

	// Create email sender
	sender := smtp.NewEmailSender(config)

	// Example 1: Plain text email
	plainMessage := smtp.EmailMessage{
		To:        []string{"ryansat46@gmail.com"},
		Subject:   "Hello from Go SMTP",
		PlainBody: "This is a test email sent from Go using SMTP.",
	}

	// Example 2: HTML email with attachment
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; }
			.header { color: #0066cc; font-size: 24px; }
		</style>
	</head>
	<body>
		<div class="header">Hello from Go SMTP!</div>
		<p>This is a <strong>HTML</strong> email sent from Go.</p>
	</body>
	</html>
	`

	htmlMessage := smtp.EmailMessage{
		To:        []string{"ryansat46@gmail.com"},
		Cc:        []string{"cc-recipient@example.com"},
		Bcc:       []string{"bcc-recipient@example.com"},
		Subject:   "HTML Email with Attachment",
		PlainBody: "This is the plain text version for email clients that don't support HTML.",
		HTMLBody:  htmlBody,
		Attachments: []smtp.Attachment{
			smtp.CreateAttachmentFromBytes("test.txt", "text/plain", []byte("This is a test attachment")),
		},
	}

	// Print example 1 details
	fmt.Println("Example 1: Plain Text Email")
	fmt.Printf("To: %v\n", plainMessage.To)
	fmt.Printf("Subject: %s\n", plainMessage.Subject)
	fmt.Printf("Body: %s\n\n", plainMessage.PlainBody)

	// Print example 2 details
	fmt.Println("Example 2: HTML Email with Attachment")
	fmt.Printf("To: %v\n", htmlMessage.To)
	fmt.Printf("Cc: %v\n", htmlMessage.Cc)
	fmt.Printf("Subject: %s\n", htmlMessage.Subject)
	fmt.Printf("Attachments: %d\n", len(htmlMessage.Attachments))
	fmt.Println("HTML Body: [HTML content not displayed]")

	// Actually send the email
	fmt.Println("Sending email to", plainMessage.To[0])
	fmt.Println("This may take a moment...")
	
	// Add a timeout context for the email sending
	err := sender.SendEmail(plainMessage)
	if err != nil {
		fmt.Printf("Failed to send plain email: %v\n", err)
		fmt.Println("This could be due to:")
		fmt.Println("- Incorrect SMTP credentials")
		fmt.Println("- Network connectivity issues")
		fmt.Println("- SMTP server requiring TLS/SSL")
	} else {
		fmt.Println("Email sent successfully!")
	}

	// Uncomment to send HTML email as well
	// err = sender.SendEmail(htmlMessage)
	// if err != nil {
	// 	fmt.Printf("Failed to send HTML email: %v\n", err)
	// }

	fmt.Println("\nTo send actual emails:")
	fmt.Println("1. Update the configuration with your SMTP server details")
	fmt.Println("2. Uncomment the sending code")
	fmt.Println("3. Run the program")

	// Prevent unused variable warning
	_ = sender
}