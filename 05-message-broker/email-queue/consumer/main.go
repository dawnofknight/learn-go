package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

const (
	headerAttempts = "x-attempts"
	maxAttempts    = 5
)

func loadEnv() {
	// Try to load .env from current directory first, then parent directory
	envPaths := []string{".env", "../.env"}
	var file *os.File
	var err error
	
	for _, path := range envPaths {
		file, err = os.Open(path)
		if err == nil {
			break
		}
	}
	
	if err != nil {
		// .env file not found, use system environment variables
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}

func mustEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	loadEnv() // Load environment variables from .env file
	amqpURL := mustEnv("AMQP_URL", "amqp://guest:guest@localhost:5672/")
	smtpHost := mustEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := mustEnv("SMTP_PORT", "587")
	smtpUser := mustEnv("SMTP_USER", "")
	smtpPass := mustEnv("SMTP_PASS", "")
	from := mustEnv("SMTP_FROM", smtpUser)

	conn, err := amqp.Dial(amqpURL)
	must(err, "dial")
	defer conn.Close()

	ch, err := conn.Channel()
	must(err, "channel")
	defer ch.Close()

	declareTopology(ch)
	must(ch.Qos(10, 0, false), "qos")

	msgs, err := ch.Consume("emails.primary", "", false, false, false, false, nil)
	must(err, "consume")

	log.Println("Worker running...")
	for d := range msgs {
		attempts := getAttempts(d.Headers)

		var job EmailJob
		if err := json.Unmarshal(d.Body, &job); err != nil {
			log.Printf("bad payload: %v", err)
			deadLetter(ch, d, attempts+1)
			_ = d.Ack(false)
			continue
		}

		if err := sendSMTP(smtpHost, smtpPort, smtpUser, smtpPass, from, job); err != nil {
			log.Printf("send error (attempt %d): %v", attempts+1, err)
			if attempts+1 >= maxAttempts {
				deadLetter(ch, d, attempts+1)
			} else {
				retry(ch, d, attempts+1)
			}
			_ = d.Ack(false) // we republished
			continue
		}

		log.Printf("email sent to %s", job.To)
		_ = d.Ack(false)
	}
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

func getAttempts(h amqp.Table) int {
	if h == nil {
		return 0
	}
	if v, ok := h[headerAttempts]; ok {
		switch t := v.(type) {
		case int32:
			return int(t)
		case int64:
			return int(t)
		case int:
			return t
		case string:
			if n, err := strconv.Atoi(t); err == nil {
				return n
			}
		}
	}
	return 0
}

func retry(ch *amqp.Channel, d amqp.Delivery, attempts int) {
	headers := d.Headers
	if headers == nil {
		headers = amqp.Table{}
	}
	headers[headerAttempts] = int32(attempts)

	_ = ch.PublishWithContext(context.Background(), "emails.dlx", "retry", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
		Timestamp:    time.Now(),
	})
}

func deadLetter(ch *amqp.Channel, d amqp.Delivery, attempts int) {
	headers := d.Headers
	if headers == nil {
		headers = amqp.Table{}
	}
	headers[headerAttempts] = int32(attempts)

	_ = ch.PublishWithContext(context.Background(), "emails.dlx", "dead", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
		Timestamp:    time.Now(),
	})
}

func sendSMTP(host, port, user, pass, from string, job EmailJob) error {
	addr := net.JoinHostPort(host, port)

	// Create email message with sender name
	var fromHeader string
	if smtpSenderName := mustEnv("SMTP_SENDER_NAME", ""); smtpSenderName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", smtpSenderName, from)
	} else {
		fromHeader = from
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n",
		fromHeader, job.To, job.Subject, job.Body,
	))
	auth := smtp.PlainAuth("", user, pass, host)
	return smtp.SendMail(addr, auth, from, []string{job.To}, msg)
}

func must(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
