# SMTP Email Sender Example

This example demonstrates how to use the SMTP package to send emails in a real application.

## Prerequisites

- Go installed on your system
- SMTP server credentials (username and password)
- Recipient email address

## Running the Example

1. Set the required environment variables:

```bash
export SMTP_SERVER="smtp.example.com"
export SMTP_USER="your-email@example.com"
export SMTP_PASSWORD="your-password"
export SENDER_EMAIL="your-email@example.com"  # Optional, defaults to SMTP_USER
export SENDER_NAME="Your Name"  # Optional
export RECIPIENT_EMAIL="recipient@example.com"
```

2. Run the example:

```bash
go run send_email.go
```

## Security Notes

- Never hardcode SMTP credentials in your code
- Use environment variables or a secure configuration system
- For Gmail and some other providers, you may need to enable "Less secure app access" or use an App Password
- For production use, consider using a dedicated email service provider with an API

## Common SMTP Servers

- Gmail: `smtp.gmail.com:587`
- Outlook/Hotmail: `smtp-mail.outlook.com:587`
- Yahoo Mail: `smtp.mail.yahoo.com:587`
- Office 365: `smtp.office365.com:587`