Here’s a **single, copy‑pasteable README blueprint** you can use to scaffold two separate Go apps (Producer + Consumer) that send emails via SMTP using **RabbitMQ** with retries & DLQ. It includes the folder structure, Docker Compose, environment variables, queue topology, code skeletons, and run/test steps.

---

# RabbitMQ Email Queue (Go Producer + Go Consumer)

A minimal, production‑leaning blueprint for an **email queue** using:

* **Go Producer** → publishes email jobs to RabbitMQ
* **Go Consumer (Worker)** → consumes jobs and sends email via **SMTP**
* **RabbitMQ topology** with **retry** (TTL) + **DLQ** (dead‑letter queue)
* Optional **HTTP enqueue endpoint** (for your app/backoffice to enqueue emails)

---

## 📁 Project Structure

```
email-queue/
├─ docker-compose.yml
├─ README.md                # ← this file (blueprint)
├─ .env.example             # sample envs
├─ producer/
│  ├─ go.mod
│  └─ main.go               # publishes EmailJob
└─ consumer/
   ├─ go.mod
   └─ main.go               # consumes EmailJob and sends via SMTP
```

---

## 🔧 Stack & Concepts

* **RabbitMQ (AMQP)**: durable queues, routing, acknowledgments, retries, DLQ.
* **Two Exchanges**:

  * `emails` (direct): normal publish/consume
  * `emails.dlx` (direct): dead‑letter routing
* **Three Queues**:

  * `emails.primary` → main consumer queue
  * `emails.retry` → retry with TTL (e.g., 30s), dead‑letters back to `emails`/`send`
  * `emails.dlq` → poison messages after `maxAttempts`
* **Headers**:

  * `x-attempts` (int) on each message to count retries.

---

## 🐇 Docker Compose (RabbitMQ + UI)

`docker-compose.yml`:

```yaml
version: "3.8"
services:
  rabbitmq:
    image: rabbitmq:3.13-management
    container_name: rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"  # http://localhost:15672 (guest/guest)
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
```

Start:

```bash
docker compose up -d
```

UI: [http://localhost:15672](http://localhost:15672)  (user/pass: guest/guest)

---

## 🔐 Environment Variables

Create `.env` files (or export in shell). Example:

`.env.example`

```env
# Shared
AMQP_URL=amqp://guest:guest@localhost:5672/

# Consumer SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your@gmail.com
SMTP_PASS=your-app-password
SMTP_FROM=your@gmail.com
```

Copy to your local:

```bash
cp .env.example producer/.env
cp .env.example consumer/.env
```

(Adjust values per app.)

---

## 📨 Message Schema

`EmailJob` JSON body:

```json
{
  "to": "user@example.com",
  "subject": "Welcome",
  "body": "Hello from RabbitMQ + Go!"
}
```

AMQP headers:

* `x-attempts` (int) — retry counter.

---

## 🧭 Queue Topology (Flow)

```
Producer --publish(send)--> [Exchange: emails] --(key=send)--> [Queue: emails.primary]
   |
   v
Consumer pulls from emails.primary:
  - On success: ACK (message removed)
  - On failure: publish to DLX with key "retry" (emails.dlx → emails.retry) and ACK original

[Queue: emails.retry] (TTL=30s) --DLX--> [Exchange: emails] --(key=send)--> emails.primary

If attempts >= maxAttempts:
  publish to DLX with key "dead" → [Queue: emails.dlq]
```

---

## 🟢 Producer (Go) — Skeleton

`producer/main.go` (minimal; publishes one job and exits)

```go
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

	job := EmailJob{
		To:      "someone@example.com",
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
	_ = ch.QueueBind("emails.retry",   "retry", "emails.dlx", false, nil)
	_ = ch.QueueBind("emails.dlq",     "dead",  "emails.dlx", false, nil)
}

func must(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
```

Init & run:

```bash
cd producer
go mod init producer
go get github.com/rabbitmq/amqp091-go
AMQP_URL=amqp://guest:guest@localhost:5672/ go run .
```

> **Optional**: turn the producer into an HTTP service with a `POST /enqueue` that publishes jobs; same publish logic inside the handler.

---

## 🟠 Consumer (Go) — Skeleton

`consumer/main.go` (consumes, SMTP send, retry + DLQ)

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strconv"
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

func mustEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
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

	_ = ch.QueueBind("emails.primary", "send",  "emails",     false, nil)
	_ = ch.QueueBind("emails.retry",   "retry", "emails.dlx", false, nil)
	_ = ch.QueueBind("emails.dlq",     "dead",  "emails.dlx", false, nil)
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
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n",
		job.To, job.Subject, job.Body,
	))
	auth := smtp.PlainAuth("", user, pass, host)
	return smtp.SendMail(addr, auth, from, []string{job.To}, msg)
}

func must(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
```

Init & run:

```bash
cd consumer
go mod init consumer
go get github.com/rabbitmq/amqp091-go
# plus your SMTP envs:
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
export SMTP_USER=your@gmail.com
export SMTP_PASS=your-app-password
export SMTP_FROM=your@gmail.com
AMQP_URL=amqp://guest:guest@localhost:5672/ go run .
```

---

## ▶️ Run End‑to‑End

1. Start RabbitMQ:

```bash
docker compose up -d
```

2. In one terminal, run **consumer**:

```bash
cd consumer
AMQP_URL=amqp://guest:guest@localhost:5672/ \
SMTP_HOST=smtp.gmail.com SMTP_PORT=587 \
SMTP_USER=your@gmail.com SMTP_PASS=your-app-password SMTP_FROM=your@gmail.com \
go run .
```

3. In another terminal, run **producer**:

```bash
cd producer
AMQP_URL=amqp://guest:guest@localhost:5672/ go run .
```

4. Watch logs: email send success / retry / DLQ behavior.

5. Inspect queues in the RabbitMQ UI: [http://localhost:15672](http://localhost:15672)

---

## 🧪 Optional: HTTP Enqueue (Producer as API)

If you want an HTTP endpoint to create jobs:

```go
// inside producer/main.go (replace main or add a mode)
http.HandleFunc("/enqueue", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "POST only", http.StatusMethodNotAllowed); return
    }
    var job EmailJob
    if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest); return
    }
    // publish(job) using same logic as above...
})
http.ListenAndServe(":8081", nil)
```

Then:

```bash
curl -X POST localhost:8081/enqueue \
  -H "Content-Type: application/json" \
  -d '{"to":"user@example.com","subject":"Hi","body":"Hello!"}'
```

---

## 🛡️ Notes for Production

* **Idempotency**: include a `message_id` in the payload; de‑duplicate on consumer side.
* **Backoff**: use multiple retry queues with increasing TTLs (10s, 1m, 5m, …).
* **Observability**: log `x-attempts`, message IDs; export metrics.
* **Security**: secure SMTP creds, use TLS; restrict RabbitMQ users/permissions.
* **Throughput**: tune `Qos(prefetch)`, scale consumers horizontally.
* **Graceful shutdown**: handle signals, drain messages, close channels cleanly.

---

## ✅ Summary

* Producer publishes JSON `EmailJob` to `emails` exchange with routing key `send`.
* Consumer subscribes to `emails.primary`, sends email via SMTP.
* Failures are retried via `emails.retry` (TTL) up to `maxAttempts`, then sent to `emails.dlq`.
* Everything is decoupled, durable, and horizontally scalable.

---

Use this README as your **scaffold blueprint**. Drop it into your repo, fill in code from the skeletons, and you’ve got a working, retry‑capable email queue system in Go.
