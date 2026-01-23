# Customer Service Database Migrations

## Overview

This directory contains database migration files for the Customer Service.

## Migration Files

| Migration | Description |
|-----------|-------------|
| 001 | Create customers table |
| 002 | Create addresses table |
| 003 | Create cart_items and guest_cart_items tables |
| 004 | Create favorites table |
| 005 | Create payment_methods table |
| 006 | Create reviews and review_images tables |

## Running Migrations

### Using golang-migrate

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations up
migrate -path migrations -database "postgresql://user:password@localhost:5434/customer_db?sslmode=disable" up

# Run migrations down
migrate -path migrations -database "postgresql://user:password@localhost:5434/customer_db?sslmode=disable" down

# Check version
migrate -path migrations -database "postgresql://user:password@localhost:5434/customer_db?sslmode=disable" version
```

### Manual Execution

```bash
# Run all up migrations
for f in migrations/*.up.sql; do
    psql -h localhost -p 5434 -U user -d customer_db -f "$f"
done

# Run all down migrations (in reverse order)
for f in $(ls migrations/*.down.sql | sort -r); do
    psql -h localhost -p 5434 -U user -d customer_db -f "$f"
done
```

## Database Schema

### Tables

#### customers
- Primary customer profile information
- Links to user_id from Auth Service

#### addresses
- Customer shipping addresses
- Supports multiple addresses with default flag

#### cart_items
- Shopping cart for authenticated users
- Auto-expires after set duration

#### guest_cart_items
- Shopping cart for guest users (session-based)
- Auto-expires after 24 hours

#### favorites
- Customer's favorite products
- Optional restock notifications

#### payment_methods
- Stripe payment method references
- Supports multiple payment methods with default flag

#### reviews
- Product reviews by customers
- Editable for 30 days after creation

#### review_images
- Images attached to reviews
- Multiple images per review supported
