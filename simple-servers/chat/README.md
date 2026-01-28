# Chat Service

## Overview
Chat Service provides real-time messaging capabilities between customers and shop owners.

## Database
- **Instance**: postgres_chat (port 22019)
- **Database**: chat_service
- **Password**: postgres_password

## Architecture
- **Pattern**: Database per Service
- **Port**: 22109 (gRPC)
- **PostgreSQL**: localhost:22019

## Schema

### chat_rooms
- `id` (UUID, primary key)
- `customer_id` (UUID, not null)
- `shop_id` (UUID, not null)
- `status` (VARCHAR(50), default: 'active')
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### messages
- `id` (UUID, primary key)
- `chat_room_id` (UUID, foreign key)
- `sender_id` (UUID, not null)
- `sender_type` (VARCHAR(50), not null)
- `message_text` (TEXT, not null)
- `is_read` (BOOLEAN, default: false)
- `created_at` (TIMESTAMP)

## Build & Run

```bash
# Build
go build -o chat-service

# Run
./chat-service

# Environment Variables
export CHAT_DATABASE_URL="postgresql://postgres:postgres_password@localhost:22019/chat_service?sslmode=disable"
export CHAT_SERVICE_PORT="22109"
```

## Database Initialization

```bash
# Apply schema
docker exec -i go_microservice_postgres_chat_dev psql -U postgres -d chat_service < schema.sql
```
