package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func mustEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	url := mustEnv("AMQP_URL", "amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(url)
	must(err, "dial")
	defer conn.Close()

	ch, err := conn.Channel()
	must(err, "channel")
	defer ch.Close()

	declareTopology(ch)

	// Get recipient from command line argument or environment variable
	recipient := "someone@example.com" // default
	if len(os.Args) > 1 {
		recipient = os.Args[1]
	} else if envRecipient := os.Getenv("EMAIL_RECIPIENT"); envRecipient != "" {
		recipient = envRecipient
	}

	job := EmailJob{
		To:      recipient,
		Subject: "Welcome",
		Body:    "Hello from RabbitMQ + Go!",
	}
	body, _ := json.Marshal(job)

	// publisher confirm (optional but recommended)
	must(ch.Confirm(false), "publisher confirm")
	acks := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	headers := amqp.Table{"x-attempts": int32(0)}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx, "emails", "send", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
		Timestamp:    time.Now(),
	})
	must(err, "publish")

	if ack := <-acks; !ack.Ack {
		log.Fatal("publish not confirmed")
	}
	log.Println("Published 1 email job.")
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
		"x-message-ttl":             int32(30000), // 30s
	})
	_, _ = ch.QueueDeclare("emails.dlq", true, false, false, false, nil)

	_ = ch.QueueBind("emails.primary", "send", "emails", false, nil)
	_ = ch.QueueBind("emails.retry", "retry", "emails.dlx", false, nil)
	_ = ch.QueueBind("emails.dlq", "dead", "emails.dlx", false, nil)
}

func must(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
