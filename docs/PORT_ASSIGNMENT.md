# Port Assignment Guide

## Overview
All services use ports in the **22000-22299** range to avoid conflicts with standard ports.

## Port Allocation Strategy

```
22000-22009: Infrastructure Services (Elasticsearch, RabbitMQ, MinIO, MailHog)
22010-22021: PostgreSQL Databases (12 microservices)
22030-22041: Redis Caches (12 microservices)
22100-22111: Microservice gRPC Servers (12 microservices)
22200-22299: Web Applications and Frontend Services
```

## Complete Port Mapping

### Infrastructure Services (22000-22009)

| Service | Port | Purpose | Container Name |
|---------|------|---------|----------------|
| Elasticsearch HTTP | 22000 | Search and indexing | `go_microservice_elasticsearch_dev` |
| Elasticsearch Transport | 22001 | Inter-node communication | `go_microservice_elasticsearch_dev` |
| RabbitMQ AMQP | 22002 | Message broker | `go_microservice_rabbitmq_dev` |
| RabbitMQ Management UI | 22003 | Web management interface | `go_microservice_rabbitmq_dev` |
| MinIO API | 22004 | Object storage API | `go_microservice_minio_dev` |
| MinIO Console | 22005 | Web console | `go_microservice_minio_dev` |
| MailHog SMTP | 22006 | Email testing SMTP | `go_microservice_mailhog_dev` |
| MailHog UI | 22007 | Email testing web UI | `go_microservice_mailhog_dev` |

### PostgreSQL Databases (22010-22021)

| Microservice | Port | Database Name | Container Name |
|--------------|------|---------------|----------------|
| Auth | 22010 | `auth_service` | `go_microservice_postgres_auth_dev` |
| Shop | 22011 | `shop_service` | `go_microservice_postgres_shop_dev` |
| Customer | 22012 | `customer_service` | `go_microservice_postgres_customer_dev` |
| Inventory | 22013 | `inventory_service` | `go_microservice_postgres_inventory_dev` |
| Order | 22014 | `order_service` | `go_microservice_postgres_order_dev` |
| Payment | 22015 | `payment_service` | `go_microservice_postgres_payment_dev` |
| Shipping | 22016 | `shipping_service` | `go_microservice_postgres_shipping_dev` |
| Notification | 22017 | `notification_service` | `go_microservice_postgres_notification_dev` |
| Review | 22018 | `review_service` | `go_microservice_postgres_review_dev` |
| Chat | 22019 | `chat_service` | `go_microservice_postgres_chat_dev` |
| Search | 22020 | `search_service` | `go_microservice_postgres_search_dev` |
| Admin | 22021 | `admin_service` | `go_microservice_postgres_admin_dev` |

### Redis Caches (22030-22041)

| Microservice | Port | Container Name |
|--------------|------|----------------|
| Auth | 22030 | `go_microservice_redis_auth_dev` |
| Shop | 22031 | `go_microservice_redis_shop_dev` |
| Customer | 22032 | `go_microservice_redis_customer_dev` |
| Inventory | 22033 | `go_microservice_redis_inventory_dev` |
| Order | 22034 | `go_microservice_redis_order_dev` |
| Payment | 22035 | `go_microservice_redis_payment_dev` |
| Shipping | 22036 | `go_microservice_redis_shipping_dev` |
| Notification | 22037 | `go_microservice_redis_notification_dev` |
| Review | 22038 | `go_microservice_redis_review_dev` |
| Chat | 22039 | `go_microservice_redis_chat_dev` |
| Search | 22040 | `go_microservice_redis_search_dev` |
| Admin | 22041 | `go_microservice_redis_admin_dev` |

### Microservice gRPC Servers (22100-22111)

| Service | Port | Binary Location | Proto Definition |
|---------|------|-----------------|------------------|
| Auth | 22100 | `microservices/auth/auth-server` | `proto/auth_service/v1/auth-service.proto` |
| Shop | 22101 | `microservices/shop/shop-server` | `proto/shop_service/v1/shop-service.proto` |
| Customer | 22102 | TBD | TBD |
| Inventory | 22103 | TBD | TBD |
| Order | 22104 | TBD | TBD |
| Payment | 22105 | TBD | TBD |
| Shipping | 22106 | TBD | TBD |
| Notification | 22107 | TBD | TBD |
| Review | 22108 | TBD | TBD |
| Chat | 22109 | TBD | TBD |
| Search | 22110 | TBD | TBD |
| Admin | 22111 | TBD | TBD |

### Web Applications (22200-22299)

| Application | Port | Framework | Location |
|-------------|------|-----------|----------|
| Shop Mall Web | 22200 | Phoenix LiveView | `web/shop_mall_web` |

## Quick Access URLs

### Customer Interfaces
- Customer Login/Register: http://localhost:22200/auth
- Password Reset Request: http://localhost:22200/auth/password-reset
- Password Reset Confirm: http://localhost:22200/auth/password-reset/confirm/{token}

### Shop Owner Interfaces
- Owner Login/Register: http://localhost:22200/owner/auth
- Owner Dashboard: http://localhost:22200/owner/dashboard

### Infrastructure Web UIs
- RabbitMQ Management: http://localhost:22003 (guest/guest)
- MinIO Console: http://localhost:22005 (minioadmin/minioadmin)
- MailHog UI: http://localhost:22007

## Environment Variables

### Phoenix Application (.env)
```bash
# Auth Service
AUTH_SERVICE_HOST=localhost
AUTH_SERVICE_PORT=22100

# Shop Service
SHOP_SERVICE_HOST=localhost
SHOP_SERVICE_PORT=22101

# Phoenix Server
PORT=22200
```

### Docker Infrastructure (.env)
```bash
# Elasticsearch
ELASTICSEARCH_PORT=22000
ELASTICSEARCH_TRANSPORT_PORT=22001

# RabbitMQ
RABBITMQ_PORT=22002
RABBITMQ_MANAGEMENT_PORT=22003

# MinIO
MINIO_PORT=22004
MINIO_CONSOLE_PORT=22005

# MailHog
MAILHOG_SMTP_PORT=22006
MAILHOG_UI_PORT=22007

# PostgreSQL (12 services: 22010-22021)
POSTGRES_AUTH_PORT=22010
POSTGRES_SHOP_PORT=22011
# ... (see .env.example for complete list)

# Redis (12 services: 22030-22041)
REDIS_AUTH_PORT=22030
REDIS_SHOP_PORT=22031
# ... (see .env.example for complete list)
```

## Troubleshooting

### Check if Port is in Use
```bash
# Check specific port
lsof -i :22100

# Check all 22xxx ports
lsof -i :22000-22299 | grep LISTEN
```

### Kill Process Using a Port
```bash
# Find PID
lsof -i :22100 | grep LISTEN

# Kill process
kill <PID>

# Force kill if needed
kill -9 <PID>
```

### Verify gRPC Service
```bash
# List services
grpcurl -plaintext localhost:22100 list

# List methods for a service
grpcurl -plaintext localhost:22100 list auth_service.v1.AuthService

# Describe a method
grpcurl -plaintext localhost:22100 describe auth_service.v1.AuthService.Register
```

### Test Database Connection
```bash
# PostgreSQL
docker exec -it go_microservice_postgres_auth_dev psql -U postgres -d auth_service

# Redis
docker exec -it go_microservice_redis_auth_dev redis-cli
```

## Port Conflict Resolution

If you encounter "port already in use" errors:

1. **Check what's using the port:**
   ```bash
   lsof -i :22XXX
   ```

2. **Stop Docker containers:**
   ```bash
   cd infrastructure/docker
   docker-compose down
   ```

3. **Kill rogue processes:**
   ```bash
   kill <PID>
   ```

4. **Restart services in order:**
   ```bash
   # 1. Infrastructure
   docker-compose up -d

   # 2. Microservices
   cd ../../microservices/auth
   ./auth-server

   # 3. Web application
   cd ../../web/shop_mall_web
   mix phx.server
   ```

## Notes

- All services use **insecure (non-TLS) gRPC** in development
- Production should use TLS certificates
- Port range 22000-22299 is chosen to avoid conflicts with:
  - Standard services (80, 443, 3000, 4000, etc.)
  - Other development tools
  - System services

## Related Documents

- [URL Access Guide](./URL_ACCESS_GUIDE.md) - Complete list of all accessible URLs
- [Elasticsearch Setup](./ELASTICSEARCH_SETUP.md) - Japanese search configuration
- [Project Status](./PROJECT_STATUS.md) - Current implementation status
