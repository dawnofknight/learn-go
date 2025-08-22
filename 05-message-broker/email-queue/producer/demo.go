package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func sendTestEmail(recipient string) {
	if recipient == "" {
		fmt.Println("Usage: go run demo.go <recipient-email>")
		fmt.Println("Example: go run demo.go user@example.com")
		return
	}

	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare queue
	_, err = ch.QueueDeclare(
		"emails.primary", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Create test email job
	emailJob := EmailJob{
		To:      recipient,
		Subject: "Test Email from RabbitMQ Queue",
		Body:    fmt.Sprintf("Hello! This is a test email sent via RabbitMQ at %s\n\nThis email was processed by our email queue system using Brevo SMTP.", time.Now().Format("2006-01-02 15:04:05")),
	}

	// Marshal to JSON
	body, err := json.Marshal(emailJob)
	if err != nil {
		log.Fatalf("Failed to marshal email job: %v", err)
	}

	// Publish message
	err = ch.Publish(
		"",               // exchange
		"emails.primary", // routing key
		false,            // mandatory
		false,            // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	fmt.Printf("✅ Email job sent to queue successfully!\n")
	fmt.Printf("📧 Recipient: %s\n", recipient)
	fmt.Printf("📝 Subject: %s\n", emailJob.Subject)
	fmt.Printf("\n💡 Make sure the consumer is running to process this email.\n")
	fmt.Printf("   Run: cd ../consumer && go run main.go\n")
}

func init() {
	if len(os.Args) >= 2 && os.Args[1] == "demo" {
		recipient := ""
		if len(os.Args) >= 3 {
			recipient = os.Args[2]
		}
		sendTestEmail(recipient)
		os.Exit(0)
	}
}
