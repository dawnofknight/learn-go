package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test-local.go <test|check>")
		fmt.Println("  test  - Test RabbitMQ connection and queue setup")
		fmt.Println("  check - Check if RabbitMQ is accessible")
		os.Exit(1)
	}

	command := os.Args[1]
	amqpURL := "amqp://guest:guest@localhost:5672/"

	switch command {
	case "check":
		checkRabbitMQ(amqpURL)
	case "test":
		testSystem(amqpURL)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func checkRabbitMQ(amqpURL string) {
	fmt.Println("Checking RabbitMQ connection...")

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		fmt.Printf("❌ Failed to connect to RabbitMQ: %v\n", err)
		fmt.Println("\n💡 To start RabbitMQ:")
		fmt.Println("   docker-compose up -d")
		fmt.Println("   or install RabbitMQ locally")
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("❌ Failed to open channel: %v\n", err)
		return
	}
	defer ch.Close()

	fmt.Println("✅ RabbitMQ connection successful!")
	fmt.Println("📊 Management UI: http://localhost:15672 (guest/guest)")
}

func testSystem(amqpURL string) {
	fmt.Println("Testing email queue system...")

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		fmt.Printf("❌ Failed to connect to RabbitMQ: %v\n", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("❌ Failed to open channel: %v\n", err)
		return
	}
	defer ch.Close()

	// Declare topology
	fmt.Println("🔧 Setting up queues and exchanges...")
	declareTopology(ch)

	// Test message
	testEmail := EmailJob{
		To:      "ryansa46@gmail.com",
		Subject: "Test Email from Queue System",
		Body:    fmt.Sprintf("This is a test email sent at %s", time.Now().Format(time.RFC3339)),
	}

	body, err := json.Marshal(testEmail)
	if err != nil {
		fmt.Printf("❌ Failed to marshal email: %v\n", err)
		return
	}

	// Publish message
	fmt.Println("📤 Publishing test email...")
	err = ch.Publish(
		"emails",
		"send",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		fmt.Printf("❌ Failed to publish message: %v\n", err)
		return
	}

	fmt.Println("✅ Test email published successfully!")
	fmt.Println("\n📋 Next steps:")
	fmt.Println("   1. Configure SMTP settings in .env file")
	fmt.Println("   2. Run consumer: cd consumer && go run main.go")
	fmt.Println("   3. Check RabbitMQ management UI for message processing")

	// Check queue status
	fmt.Println("\n📊 Queue Status:")
	checkQueue(ch, "emails.primary")
	checkQueue(ch, "emails.retry")
	checkQueue(ch, "emails.dlq")
}

func declareTopology(ch *amqp.Channel) {
	_ = ch.ExchangeDeclare("emails", "direct", true, false, false, false, nil)
	_ = ch.ExchangeDeclare("emails.dlx", "direct", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("emails.primary", true, false, false, false, amqp.Table{
		"x-dead-letter-exchange": "emails.dlx",
	})
	_, _ = ch.QueueDeclare("emails.retry", true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    "emails",
		"x-dead-letter-routing-key": "send",
		"x-message-ttl":             int32(30000),
	})
	_, _ = ch.QueueDeclare("emails.dlq", true, false, false, false, nil)

	_ = ch.QueueBind("emails.primary", "send", "emails", false, nil)
	_ = ch.QueueBind("emails.retry", "retry", "emails.dlx", false, nil)
	_ = ch.QueueBind("emails.dlq", "dead", "emails.dlx", false, nil)
}

func checkQueue(ch *amqp.Channel, queueName string) {
	queue, err := ch.QueueInspect(queueName)
	if err != nil {
		fmt.Printf("   ❌ %s: error - %v\n", queueName, err)
		return
	}
	fmt.Printf("   📦 %s: %d messages\n", queueName, queue.Messages)
}
