# Inventory Service

## Overview
The Inventory Service manages product inventory and reservations for the microservice architecture.

## Database Configuration
- **Host**: localhost (127.0.0.1)
- **Port**: 22013
- **Database**: inventory_service
- **User**: postgres
- **Password**: postgres_password

## Service Configuration
- **Service Port**: 22103
- **Protocol**: gRPC

## Database Schema

### Tables

#### inventories
```sql
CREATE TABLE inventories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    shop_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
- `idx_inventories_product_id` on product_id
- `idx_inventories_shop_id` on shop_id

#### reservations
```sql
CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
    order_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'reserved',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
- `idx_reservations_inventory_id` on inventory_id
- `idx_reservations_order_id` on order_id

**Foreign Keys:**
- `inventory_id` references `inventories(id)` ON DELETE CASCADE

## Build & Run

### Build
```bash
go build -o inventory-service main.go
```

### Run
```bash
./inventory-service
```

### Environment Variables
- `INVENTORY_DATABASE_URL`: PostgreSQL connection string (default: postgresql://postgres:postgres_password@127.0.0.1:22013/inventory_service?sslmode=disable)
- `INVENTORY_SERVICE_PORT`: Service port (default: 22103)

## Verification

### Check Service Status
```bash
ps aux | grep inventory-service
```

### Check Database Connection
```bash
lsof -i :22013 | grep inventory
```

### Check Service Listening
```bash
lsof -i :22103
```

### Check Database Tables
```bash
docker exec -i go_microservice_postgres_inventory_dev psql -U postgres -d inventory_service -c "\dt"
```

## Architecture Notes
- Part of Database per Service architecture
- Each microservice has its own dedicated PostgreSQL instance
- Inventory Service has dedicated PostgreSQL on port 22013
- Uses gRPC for inter-service communication
