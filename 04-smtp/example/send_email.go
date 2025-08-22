package main

// Import the local smtp package
// For a real application, you would use the proper import path after setting up your module

// func main() {
// 	// Get SMTP credentials from environment variables for security
// 	smtpServer := getEnv("SMTP_SERVER", "smtp.example.com")
// 	smtpPort := 587 // Default port for TLS
// 	smtpUser := getEnv("SMTP_USER", "")
// 	smtpPass := getEnv("SMTP_PASSWORD", "")
// 	senderEmail := getEnv("SENDER_EMAIL", smtpUser)
// 	senderName := getEnv("SENDER_NAME", "Go SMTP Example")

// 	// Check if credentials are provided
// 	if smtpUser == "" || smtpPass == "" {
// 		log.Fatal("SMTP credentials not provided. Please set SMTP_USER and SMTP_PASSWORD environment variables.")
// 	}

// 	// Create email configuration
// 	config := smtp.EmailConfig{
// 		SMTPServer:   smtpServer,
// 		SMTPPort:     smtpPort,
// 		SMTPUsername: smtpUser,
// 		SMTPPassword: smtpPass,
// 		SenderEmail:  senderEmail,
// 		SenderName:   senderName,
// 	}

// 	// Create email sender
// 	sender := smtp.NewEmailSender(config)

// 	// Create a simple email message
// 	message := smtp.EmailMessage{
// 		To:        []string{getEnv("RECIPIENT_EMAIL", "")},
// 		Subject:   "Test Email from Go SMTP Example",
// 		PlainBody: "This is a test email sent from the Go SMTP example application.",
// 		HTMLBody:  "<h1>Test Email</h1><p>This is a <strong>test email</strong> sent from the Go SMTP example application.</p>",
// 	}

// 	// Check if recipient is provided
// 	if len(message.To) == 0 || message.To[0] == "" {
// 		log.Fatal("Recipient email not provided. Please set RECIPIENT_EMAIL environment variable.")
// 	}

// 	// Send the email
// 	fmt.Println("Sending email to", message.To[0])
// 	err := sender.SendEmail(message)
// 	if err != nil {
// 		log.Fatalf("Failed to send email: %v", err)
// 	}

// 	fmt.Println("Email sent successfully!")
// }

// // getEnv gets an environment variable or returns a default value
// func getEnv(key, defaultValue string) string {
// 	value := os.Getenv(key)
// 	if value == "" {
// 		return defaultValue
// 	}
// 	return value
// }
