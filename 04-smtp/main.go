// Package smtp provides SMTP email functionality
package smtp

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EmailConfig holds the configuration for SMTP email sending
type EmailConfig struct {
	SMTPServer         string
	SMTPPort           int
	SMTPUsername       string
	SMTPPassword       string
	SenderEmail        string
	SenderName         string
	InsecureSkipVerify bool // Skip TLS certificate verification (for testing only)
	DebugMode          bool // Enable debug logging
	AuthMethod         string // Authentication method: "plain", "login", or "cram-md5"
}

// EmailMessage represents an email message to be sent
type EmailMessage struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	PlainBody   string
	HTMLBody    string
	Attachments []Attachment
}

// Attachment represents a file attachment for an email
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// EmailSender handles sending emails via SMTP
type EmailSender struct {
	Config EmailConfig
}

// loginAuth is a custom implementation of smtp.Auth for LOGIN authentication
type loginAuth struct {
	username, password string
}

// Start implements smtp.Auth.Start method
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next implements smtp.Auth.Next method
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:", "334 VXNlcm5hbWU6": // Base64 encoded "Username:"
			return []byte(a.username), nil
		case "Password:", "334 UGFzc3dvcmQ6": // Base64 encoded "Password:"
			return []byte(a.password), nil
		default:
			// Try to decode the challenge
			decodedChallenge, err := base64.StdEncoding.DecodeString(string(fromServer[4:]))
			if err == nil {
				if strings.Contains(strings.ToLower(string(decodedChallenge)), "username") {
					return []byte(a.username), nil
				} else if strings.Contains(strings.ToLower(string(decodedChallenge)), "password") {
					return []byte(a.password), nil
				}
			}
			return nil, fmt.Errorf("unknown challenge: %s", fromServer)
		}
	}
	return nil, nil
}

// NewEmailSender creates a new email sender with the given configuration
func NewEmailSender(config EmailConfig) *EmailSender {
	return &EmailSender{Config: config}
}

// SendEmail sends an email using the configured SMTP server
func (s *EmailSender) SendEmail(message EmailMessage) error {
	// Validate required fields
	if len(message.To) == 0 {
		return fmt.Errorf("recipient email address is required")
	}

	if message.Subject == "" {
		return fmt.Errorf("email subject is required")
	}

	if message.PlainBody == "" && message.HTMLBody == "" {
		return fmt.Errorf("email body (plain or HTML) is required")
	}

	// Debug logging
	if s.Config.DebugMode {
		fmt.Println("[DEBUG] Starting email send process")
		fmt.Printf("[DEBUG] SMTP Server: %s:%d\n", s.Config.SMTPServer, s.Config.SMTPPort)
		fmt.Printf("[DEBUG] Username: %s\n", s.Config.SMTPUsername)
		fmt.Printf("[DEBUG] From: %s <%s>\n", s.Config.SenderName, s.Config.SenderEmail)
		fmt.Printf("[DEBUG] To: %v\n", message.To)
		fmt.Printf("[DEBUG] Subject: %s\n", message.Subject)
		fmt.Printf("[DEBUG] InsecureSkipVerify: %v\n", s.Config.InsecureSkipVerify)
	}

	// Create email content
	email := s.buildEmail(message)

	// Prepare recipient list
	recipients := append(append(message.To, message.Cc...), message.Bcc...)
	
	// Format SMTP server address
	smtpAddr := fmt.Sprintf("%s:%d", s.Config.SMTPServer, s.Config.SMTPPort)
	
	// Set up authentication based on the specified method
	var auth smtp.Auth
	
	switch s.Config.AuthMethod {
	case "cram-md5":
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Using CRAM-MD5 authentication")
		}
		auth = smtp.CRAMMD5Auth(s.Config.SMTPUsername, s.Config.SMTPPassword)
		
	case "login":
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Using LOGIN authentication")
		}
		// Use our custom LOGIN authentication implementation
		auth = &loginAuth{
			username: s.Config.SMTPUsername,
			password: s.Config.SMTPPassword,
		}
		
	default: // "plain" or empty
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Using PLAIN authentication")
		}
		auth = smtp.PlainAuth("", s.Config.SMTPUsername, s.Config.SMTPPassword, s.Config.SMTPServer)
	}
	
	// Check if we're using a secure port (465 is typically SMTPS)
	if s.Config.SMTPPort == 465 {
		// For SMTPS (SMTP over SSL/TLS), we need to use a different approach
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Using SMTPS (SMTP over SSL/TLS) for port 465")
		}
		
		// Create TLS config
		tlsConfig := &tls.Config{
			ServerName:         s.Config.SMTPServer,
			InsecureSkipVerify: s.Config.InsecureSkipVerify,
		}
		
		if s.Config.DebugMode {
			fmt.Printf("[DEBUG] TLS Config: ServerName=%s, InsecureSkipVerify=%v\n", 
				tlsConfig.ServerName, tlsConfig.InsecureSkipVerify)
		}
		
		// Connect to the server
		if s.Config.DebugMode {
			fmt.Printf("[DEBUG] Connecting to %s...\n", smtpAddr)
		}
		
		conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
		if err != nil {
			if s.Config.DebugMode {
				fmt.Printf("[DEBUG] Connection failed: %v\n", err)
			}
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Connection established successfully")
		}
		
		defer conn.Close()
		
		// Create a new SMTP client
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Creating SMTP client...")
		}
		
		c, err := smtp.NewClient(conn, s.Config.SMTPServer)
		if err != nil {
			if s.Config.DebugMode {
				fmt.Printf("[DEBUG] Failed to create SMTP client: %v\n", err)
			}
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] SMTP client created successfully")
		}
		
		defer c.Close()
		
		// Authenticate
		if s.Config.DebugMode {
			fmt.Printf("[DEBUG] Authenticating with username: %s\n", s.Config.SMTPUsername)
			fmt.Printf("[DEBUG] Password length: %d characters\n", len(s.Config.SMTPPassword))
			// Safely display part of the password for debugging
			passLen := len(s.Config.SMTPPassword)
			if passLen > 0 {
				firstChar := string(s.Config.SMTPPassword[0])
				fmt.Printf("[DEBUG] Password first char: %s\n", firstChar)
				if passLen > 1 {
					lastChar := string(s.Config.SMTPPassword[passLen-1])
					fmt.Printf("[DEBUG] Password last char: %s\n", lastChar)
				}
			}
		}
		
		// Use standard Auth method for all authentication types
		if s.Config.DebugMode {
			fmt.Printf("[DEBUG] Authenticating with method: %s\n", s.Config.AuthMethod)
		}
		
		if err = c.Auth(auth); err != nil {
			if s.Config.DebugMode {
				fmt.Printf("[DEBUG] Authentication failed: %v\n", err)
			}
			return fmt.Errorf("SMTP authentication failed for user %s on server %s:%d: %w", 
				s.Config.SMTPUsername, s.Config.SMTPServer, s.Config.SMTPPort, err)
		}
		
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Authentication successful")
		}
		
		if s.Config.DebugMode {
			fmt.Println("[DEBUG] Authentication successful")
		}
		
		// Set the sender and recipients
		if s.Config.DebugMode {
			fmt.Printf("[DEBUG] Setting sender: %s\n", s.Config.SenderEmail)
		}
		
		if err = c.Mail(s.Config.SenderEmail); err != nil {
			if s.Config.DebugMode {
				fmt.Printf("[DEBUG] Failed to set sender: %v\n", err)
			}
			return fmt.Errorf("failed to set sender: %w", err)
		}
		
		for _, recipient := range recipients {
			if s.Config.DebugMode {
				fmt.Printf("[DEBUG] Setting recipient: %s\n", recipient)
			}
			
			if err = c.Rcpt(recipient); err != nil {
				if s.Config.DebugMode {
					fmt.Printf("[DEBUG] Failed to set recipient: %v\n", err)
				}
				return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
			}
		}
		
		// Send the email body
		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("failed to open data writer: %w", err)
		}
		
		_, err = w.Write([]byte(email))
		if err != nil {
			return fmt.Errorf("failed to write email data: %w", err)
		}
		
		err = w.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}
		
		// Send the QUIT command and close the connection
		err = c.Quit()
		if err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	} else if s.Config.SMTPPort == 587 {
		// For SMTP with STARTTLS (port 587)
		// Connect to the server
		c, err := smtp.Dial(smtpAddr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer c.Close()
		
		// Start TLS
		tlsConfig := &tls.Config{
			ServerName:         s.Config.SMTPServer,
			InsecureSkipVerify: s.Config.InsecureSkipVerify,
		}
		if err = c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
		
		// Authenticate
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed for user %s on server %s:%d: %w", 
				s.Config.SMTPUsername, s.Config.SMTPServer, s.Config.SMTPPort, err)
		}
		
		// Set the sender and recipients
		if err = c.Mail(s.Config.SenderEmail); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}
		
		for _, recipient := range recipients {
			if err = c.Rcpt(recipient); err != nil {
				return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
			}
		}
		
		// Send the email body
		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("failed to open data writer: %w", err)
		}
		
		_, err = w.Write([]byte(email))
		if err != nil {
			return fmt.Errorf("failed to write email data: %w", err)
		}
		
		err = w.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}
		
		// Send the QUIT command and close the connection
		err = c.Quit()
		if err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	} else {
		// For standard SMTP without encryption
		err := smtp.SendMail(smtpAddr, auth, s.Config.SenderEmail, recipients, []byte(email))
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}

// buildEmail constructs the full email content including headers and body
func (s *EmailSender) buildEmail(message EmailMessage) string {
	// Generate a boundary for multipart messages
	boundary := "==_GoEmailBoundary_" + time.Now().Format("20060102150405") + "_=="

	// Build email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.Config.SenderName, s.Config.SenderEmail)
	headers["To"] = strings.Join(message.To, ", ")
	if len(message.Cc) > 0 {
		headers["Cc"] = strings.Join(message.Cc, ", ")
	}
	headers["Subject"] = message.Subject
	headers["MIME-Version"] = "1.0"

	// Determine content type based on message content
	hasAttachments := len(message.Attachments) > 0
	hasHTML := message.HTMLBody != ""

	if hasAttachments || hasHTML {
		// Multipart email
		headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary)
	} else {
		// Simple plain text email
		headers["Content-Type"] = "text/plain; charset=UTF-8"
	}

	// Build the email content
	var emailContent strings.Builder

	// Add headers
	for key, value := range headers {
		emailContent.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	emailContent.WriteString("\r\n")

	// For simple plain text emails without attachments
	if !hasAttachments && !hasHTML {
		emailContent.WriteString(message.PlainBody)
		return emailContent.String()
	}

	// For multipart emails
	// Add plain text part if available
	if message.PlainBody != "" {
		emailContent.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		emailContent.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		emailContent.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		emailContent.WriteString(message.PlainBody)
		emailContent.WriteString("\r\n")
	}

	// Add HTML part if available
	if hasHTML {
		emailContent.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		emailContent.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		emailContent.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		emailContent.WriteString(message.HTMLBody)
		emailContent.WriteString("\r\n")
	}

	// Add attachments
	for _, attachment := range message.Attachments {
		emailContent.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		emailContent.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", 
			attachment.ContentType, attachment.Filename))
		emailContent.WriteString("Content-Transfer-Encoding: base64\r\n")
		emailContent.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", 
			attachment.Filename))

		// Encode attachment data as base64
		encoded := base64.StdEncoding.EncodeToString(attachment.Data)
		
		// Split the base64 data into lines of 76 characters
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			emailContent.WriteString(encoded[i:end] + "\r\n")
		}
	}

	// Close the multipart message
	emailContent.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return emailContent.String()
}

// CreateAttachmentFromFile creates an attachment from a file on disk
func CreateAttachmentFromFile(filePath string) (Attachment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Attachment{}, fmt.Errorf("failed to open attachment file: %w", err)
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		return Attachment{}, fmt.Errorf("failed to read attachment file: %w", err)
	}

	// Determine content type based on file extension
	contentType := "application/octet-stream" // Default content type
	ext := filepath.Ext(filePath)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".pdf":
		contentType = "application/pdf"
	case ".txt":
		contentType = "text/plain"
	case ".html", ".htm":
		contentType = "text/html"
	case ".doc", ".docx":
		contentType = "application/msword"
	case ".xls", ".xlsx":
		contentType = "application/vnd.ms-excel"
	}

	return Attachment{
		Filename:    filepath.Base(filePath),
		ContentType: contentType,
		Data:        data,
	}, nil
}

// CreateAttachmentFromBytes creates an attachment from byte data
func CreateAttachmentFromBytes(filename, contentType string, data []byte) Attachment {
	return Attachment{
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}
}

func main() {
	// Print header
	fmt.Println("=== Email Sending Example ===")
	fmt.Println("This example shows how to send emails using Go's SMTP package.")
	fmt.Println("To actually send emails, uncomment the sending code and provide real credentials.")

	// Example configuration
	config := EmailConfig{
		SMTPServer:   "smtp.example.com",
		SMTPPort:     587,
		SMTPUsername: "your-email@example.com",
		SMTPPassword: "your-password",
		SenderEmail:  "your-email@example.com",
		SenderName:   "Your Name",
	}

	// Create email sender - just for demonstration, not using it directly
	_ = NewEmailSender(config)

	// Example 1: Simple plain text email
	plainMessage := EmailMessage{
		To:        []string{"recipient@example.com"},
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
			.content { margin: 20px 0; }
			.footer { color: #666; font-size: 12px; }
		</style>
	</head>
	<body>
		<div class="header">Hello from Go SMTP!</div>
		<div class="content">
			<p>This is a <strong>HTML</strong> email sent from Go using SMTP.</p>
			<p>It demonstrates how to send formatted emails with attachments.</p>
		</div>
		<div class="footer">This is an example email. Please do not reply.</div>
	</body>
	</html>
	`

	// Create an in-memory attachment example
	exampleData := []byte("This is example attachment data.")
	attachment := CreateAttachmentFromBytes("example.txt", "text/plain", exampleData)

	htmlMessage := EmailMessage{
		To:          []string{"recipient@example.com"},
		Cc:          []string{"cc-recipient@example.com"},
		Subject:     "HTML Email with Attachment",
		PlainBody:   "This is the plain text version for email clients that don't support HTML.",
		HTMLBody:    htmlBody,
		Attachments: []Attachment{attachment},
	}

	// Uncomment to send the emails (replace with real credentials first)
	/*
	// Send plain text email
	err := sender.SendEmail(plainMessage)
	if err != nil {
		log.Printf("Failed to send plain email: %v", err)
	} else {
		log.Println("Plain email sent successfully!")
	}

	// Send HTML email with attachment
	err = sender.SendEmail(htmlMessage)
	if err != nil {
		log.Printf("Failed to send HTML email: %v", err)
	} else {
		log.Println("HTML email with attachment sent successfully!")
	}
	*/

	// Print example 1 details
	fmt.Println("\nExample 1: Plain Text Email")
	fmt.Printf("To: %s\n", strings.Join(plainMessage.To, ", "))
	fmt.Printf("Subject: %s\n", plainMessage.Subject)
	fmt.Printf("Body: %s\n", plainMessage.PlainBody)

	// Print example 2 details
	fmt.Println("\nExample 2: HTML Email with Attachment")
	fmt.Printf("To: %s\n", strings.Join(htmlMessage.To, ", "))
	if len(htmlMessage.Cc) > 0 {
		fmt.Printf("Cc: %s\n", strings.Join(htmlMessage.Cc, ", "))
	}
	fmt.Printf("Subject: %s\n", htmlMessage.Subject)
	fmt.Printf("Attachments: %d\n", len(htmlMessage.Attachments))
	fmt.Println("HTML Body: [HTML content not displayed]")
}