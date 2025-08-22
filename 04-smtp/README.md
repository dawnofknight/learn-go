# Go SMTP Email Sender

This package provides a simple yet powerful implementation for sending emails using Go's built-in SMTP package. It supports plain text emails, HTML emails, and file attachments.

## Features

- Send plain text emails
- Send HTML formatted emails
- Add file attachments
- Support for CC and BCC recipients
- Configurable SMTP settings
- Error handling and validation

## Usage

### Basic Configuration

```go
// Create email configuration
config := EmailConfig{
    SMTPServer:   "smtp.example.com",
    SMTPPort:     587,
    SMTPUsername: "your-email@example.com",
    SMTPPassword: "your-password",
    SenderEmail:  "your-email@example.com",
    SenderName:   "Your Name",
}

// Create email sender
sender := NewEmailSender(config)
```

### Sending a Simple Plain Text Email

```go
// Create a plain text email
plainMessage := EmailMessage{
    To:        []string{"recipient@example.com"},
    Subject:   "Hello from Go SMTP",
    PlainBody: "This is a test email sent from Go using SMTP.",
}

// Send the email
err := sender.SendEmail(plainMessage)
if err != nil {
    log.Printf("Failed to send email: %v", err)
} else {
    log.Println("Email sent successfully!")
}
```

### Sending an HTML Email with Attachment

```go
// Create HTML content
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

// Create an attachment from a file
attachment, err := CreateAttachmentFromFile("document.pdf")
if err != nil {
    log.Fatalf("Failed to create attachment: %v", err)
}

// Create the email message
htmlMessage := EmailMessage{
    To:          []string{"recipient@example.com"},
    Cc:          []string{"cc-recipient@example.com"},
    Subject:     "HTML Email with Attachment",
    PlainBody:   "This is the plain text version for email clients that don't support HTML.",
    HTMLBody:    htmlBody,
    Attachments: []Attachment{attachment},
}

// Send the email
err = sender.SendEmail(htmlMessage)
if err != nil {
    log.Printf("Failed to send email: %v", err)
} else {
    log.Println("Email sent successfully!")
}
```

### Creating Attachments

From a file:
```go
attachment, err := CreateAttachmentFromFile("path/to/file.pdf")
if err != nil {
    log.Fatalf("Failed to create attachment: %v", err)
}
```

From byte data:
```go
data := []byte("This is the content of the attachment")
attachment := CreateAttachmentFromBytes("filename.txt", "text/plain", data)
```

## Common SMTP Servers

- Gmail: `smtp.gmail.com:587`
- Outlook/Hotmail: `smtp-mail.outlook.com:587`
- Yahoo Mail: `smtp.mail.yahoo.com:587`
- Office 365: `smtp.office365.com:587`

## Security Notes

1. Never hardcode email credentials in your code
2. Consider using environment variables or a secure configuration system
3. For Gmail and some other providers, you may need to enable "Less secure app access" or use an App Password
4. For production use, consider using a dedicated email service provider with an API

## Error Handling

The `SendEmail` function returns an error if:
- Required fields are missing (recipient, subject, body)
- SMTP authentication fails
- Connection to the SMTP server fails
- Any other error occurs during the sending process

Always check the returned error to ensure emails are sent successfully.