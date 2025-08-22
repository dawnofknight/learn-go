# Email Queue System with RabbitMQ

A robust email queue system built with Go and RabbitMQ that provides reliable email delivery with retry mechanisms and dead letter handling.

## Features

- **Reliable Message Delivery**: Uses RabbitMQ with publisher confirms
- **Retry Mechanism**: Automatic retry with exponential backoff
- **Dead Letter Queue**: Failed messages are moved to DLQ after max attempts
- **SMTP Integration**: Sends emails via SMTP with configurable providers
- **Environment Configuration**: Easy configuration via environment variables

## Architecture

```
Producer → emails exchange → emails.primary queue → Consumer
                                     ↓ (on failure)
                            emails.dlx exchange
                                     ↓
                            emails.retry queue (30s TTL)
                                     ↓ (after TTL)
                            emails exchange → emails.primary queue
                                     ↓ (max attempts reached)
                            emails.dlq (dead letter queue)
```

## Quick Start

### 1. Start RabbitMQ

```bash
docker-compose up -d
```

This will start RabbitMQ with management UI available at http://localhost:15672 (guest/guest).

### 2. Configure Environment

The project comes with pre-configured SMTP settings using Brevo (formerly SendinBlue) for immediate testing. The `.env` file is already set up with working credentials:

```env
# RabbitMQ Configuration
AMQP_URL=amqp://guest:guest@localhost:5672/

# SMTP Configuration (Brevo)
SMTP_HOST=smtp-relay.brevo.com
SMTP_PORT=2525
SMTP_USER=8b0803001@smtp-brevo.com
SMTP_PASS=GIka0UMp6Vm35wYv
SMTP_FROM=hello@ezclass.io
SMTP_SENDER_NAME=EZClass
```

To use your own SMTP provider, copy the template and modify:

```bash
cp .env.example .env
# Then edit .env with your settings
```

### 3. Run the Consumer

```bash
cd consumer
go run main.go
```

### 4. Send Test Emails

**Quick Test with Demo Script:**

```bash
# Send a test email using the demo script
cd producer
go run main.go demo your-email@example.com
```

**Manual Testing:**

```bash
cd producer
go run main.go
```

The consumer will automatically process the email and send it via SMTP using the configured Brevo credentials.

**Check Email Delivery:**
- Check your inbox for the test email
- Monitor the consumer logs for delivery status
- Use RabbitMQ Management UI at http://localhost:15672 to monitor queues

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AMQP_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection URL |
| `SMTP_HOST` | `smtp.gmail.com` | SMTP server hostname |
| `SMTP_PORT` | `587` | SMTP server port |
| `SMTP_USER` | | SMTP username |
| `SMTP_PASS` | | SMTP password |
| `SMTP_FROM` | `SMTP_USER` | From email address |

### SMTP Providers

#### Gmail
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

#### Outlook/Hotmail
```env
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=your-email@outlook.com
SMTP_PASS=your-password
```

#### Custom SMTP
```env
SMTP_HOST=your-smtp-server.com
SMTP_PORT=587
SMTP_USER=your-username
SMTP_PASS=your-password
```

## Message Format

The system expects JSON messages with the following structure:

```json
{
  "to": "recipient@example.com",
  "subject": "Test Email",
  "body": "This is a test email message."
}
```

## Retry Logic

- **Max Attempts**: 5 retries per message
- **Retry Delay**: 30 seconds (configurable via TTL)
- **Dead Letter**: Messages exceeding max attempts are moved to DLQ
- **Attempt Tracking**: Uses `x-attempts` header to track retry count

## Monitoring

### RabbitMQ Management UI

Access the management interface at http://localhost:15672 to monitor:
- Queue depths
- Message rates
- Consumer status
- Dead letter queues

### Queue Overview

- `emails.primary`: Main processing queue
- `emails.retry`: Temporary queue for failed messages (30s TTL)
- `emails.dlq`: Dead letter queue for permanently failed messages

## Development

### Project Structure

```
email-queue/
├── docker-compose.yml    # RabbitMQ setup
├── .env.example         # Environment template
├── producer/
│   ├── go.mod
│   └── main.go          # Message publisher
└── consumer/
    ├── go.mod
    └── main.go          # Email processor
```

### Building

```bash
# Build producer
cd producer && go build -o ../bin/producer

# Build consumer
cd consumer && go build -o ../bin/consumer
```

### Testing

1. Start RabbitMQ: `docker-compose up -d`
2. Configure `.env` with test SMTP settings
3. Run consumer: `cd consumer && go run main.go`
4. Send test message: `cd producer && go run main.go`
5. Check logs and RabbitMQ management UI

## Troubleshooting

### Common Issues

1. **SMTP Authentication Failed**
   - Verify SMTP credentials
   - Enable "Less secure app access" for Gmail
   - Use app-specific passwords for Gmail

2. **Connection Refused**
   - Ensure RabbitMQ is running: `docker-compose ps`
   - Check AMQP_URL configuration

3. **Messages Stuck in Retry Queue**
   - Check consumer logs for errors
   - Verify SMTP configuration
   - Monitor queue depths in management UI

### Logs

The consumer provides detailed logging:
- Successful email deliveries
- Retry attempts with error details
- Dead letter notifications
- Connection status

## Production Considerations

- Use persistent volumes for RabbitMQ data
- Configure proper authentication and SSL
- Monitor queue depths and consumer health
- Set up alerting for dead letter queue growth
- Consider horizontal scaling of consumers
- Implement proper logging and metrics collection