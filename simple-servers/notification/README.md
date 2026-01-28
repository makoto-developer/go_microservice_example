# Notification Service

A microservice responsible for managing and sending notifications across multiple channels (email, SMS, push notifications).

## Features

- Template-based notification system
- Multiple notification channels (email, SMS, push)
- Notification history tracking
- Retry mechanism for failed notifications
- Dedicated PostgreSQL database

## Database Schema

### Tables

#### notification_templates
Stores reusable notification templates with variable placeholders.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| name | VARCHAR(100) | Unique template name |
| channel | VARCHAR(50) | Notification channel (email, sms, push) |
| subject | VARCHAR(255) | Email subject (null for SMS/push) |
| body_template | TEXT | Template body with variable placeholders |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

#### notifications
Tracks all notification records and their delivery status.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| template_id | UUID | Reference to notification_templates |
| recipient | VARCHAR(255) | Recipient address (email/phone) |
| channel | VARCHAR(50) | Notification channel |
| subject | VARCHAR(255) | Actual subject (for emails) |
| body | TEXT | Actual notification body |
| status | VARCHAR(50) | Status: pending, sent, failed, retrying |
| sent_at | TIMESTAMP | When notification was sent |
| error_message | TEXT | Error message if sending failed |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

### Sample Templates

Pre-configured templates include:
- `order_confirmation` - Order confirmation emails
- `payment_success` - Payment success notifications
- `shipping_update` - Shipping status updates
- `order_cancelled` - Order cancellation notices

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DATABASE_URL | postgresql://postgres:postgres_password@localhost:22017/notification_service?sslmode=disable | PostgreSQL connection string |
| SERVICE_PORT | 22106 | gRPC service port |

### Database

- **Host**: localhost:22017 (from host machine)
- **Database**: notification_service
- **User**: postgres
- **Password**: postgres_password

## Building

```bash
go mod download
go build -o notification-service
```

## Running

```bash
./notification-service
```

Expected output:
```
✅ Successfully connected to Notification database
✅ Notification Service is running on port 22106
🎯 Database per Service architecture is active!
   - Notification Service has dedicated PostgreSQL instance on port 22017
```

## Database Setup

The schema is automatically created from `schema.sql`:

```bash
docker exec -i go_microservice_postgres_notification_dev psql -U postgres -d notification_service < schema.sql
```

## Architecture

This service follows the Database per Service pattern:
- Dedicated PostgreSQL instance (port 22017)
- Dedicated Redis cache (port 22037)
- Independent data management
- No shared database access

## Integration

### Sending Notifications

The notification service will be called by other services when they need to send notifications:

- **Order Service**: Order confirmation, status updates
- **Payment Service**: Payment confirmations, failure notifications
- **Shipping Service**: Shipping updates, delivery confirmations

### Template Usage

Templates use variable placeholders in the format `{variable_name}`:

```
Subject: Order Confirmation - Order #{order_number}
Body: Dear {customer_name},
      Your order #{order_number} has been confirmed.
      Total Amount: {total_amount}
```

## Future Enhancements

- Email provider integration (SendGrid, AWS SES)
- SMS provider integration (Twilio)
- Push notification support (FCM, APNS)
- Retry logic with exponential backoff
- Rate limiting per channel
- Notification preferences management
