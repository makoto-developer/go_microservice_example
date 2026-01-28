# Admin Service - Quick Start Guide

## Prerequisites
- Docker and Docker Compose running
- Go 1.23 or later
- PostgreSQL client tools (optional)

## Quick Start (3 Steps)

### 1. Start the Database
```bash
cd ../../infrastructure/docker
docker compose up -d postgres_admin
```

### 2. Verify Setup
```bash
cd ../../simple-servers/admin
./verify.sh
```

### 3. Run the Service
```bash
./admin-service
```

Or with custom configuration:
```bash
DATABASE_URL="postgresql://postgres:postgres_password@localhost:22021/admin_service?sslmode=disable" \
SERVICE_PORT=22111 \
./admin-service
```

## What You Get

### Database (PostgreSQL on port 22021)
- **Database**: admin_service
- **Tables**:
  - `admin_users` - Administrator accounts
  - `audit_logs` - System audit logs

### Service (gRPC on port 22111)
- Database per Service architecture
- Dedicated PostgreSQL instance
- gRPC reflection enabled
- Health monitoring ready

## Verification

Run the verification script:
```bash
./verify.sh
```

Expected output:
```
✅ PostgreSQL container is running
✅ Database is accessible
✅ Required tables exist (admin_users, audit_logs)
✅ Service binary exists
✅ Service started successfully
✅ Service is listening on port 22111
```

## Database Test

Test database operations:
```bash
./test_db.sh
```

This will show:
- Table structures
- Sample data in admin_users
- Sample data in audit_logs

## Development

### Build
```bash
go build -o admin-service
```

### Test Database Schema
```bash
docker exec -i go_microservice_postgres_admin_dev psql -U postgres -d admin_service < sql/schema.sql
```

### Query Database
```bash
docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -c "
SELECT * FROM admin_users;
"
```

## Troubleshooting

### Database Connection Error
If you see password authentication errors:
```bash
docker exec go_microservice_postgres_admin_dev psql -U postgres \
  -c "ALTER USER postgres PASSWORD 'postgres_password';"
```

### Port Already in Use
Change the service port:
```bash
SERVICE_PORT=22112 ./admin-service
```

### Rebuild from Scratch
```bash
# Stop container
cd ../../infrastructure/docker
docker compose down postgres_admin

# Remove volume
docker volume rm go_microservice_postgres_admin_data_dev

# Start fresh
docker compose up -d postgres_admin

# Wait for startup
sleep 5

# Setup schema
cd ../../simple-servers/admin
./setup_db.sh
```

## Architecture

```
Admin Service (port 22111)
    ↓
PostgreSQL (port 22021)
    ├── admin_users table
    └── audit_logs table
```

## Files

```
simple-servers/admin/
├── main.go                         # Service entry point
├── go.mod                          # Go dependencies
├── admin-service                   # Compiled binary
├── sql/
│   └── schema.sql                  # Database schema
├── README.md                       # Detailed documentation
├── QUICK_START.md                  # This file
├── IMPLEMENTATION_SUMMARY.md       # Implementation details
├── setup_db.sh                     # Database setup
├── verify.sh                       # Verification script
└── test_db.sh                      # Database test script
```

## Next Steps

1. Implement gRPC service handlers
2. Add authentication middleware
3. Implement audit logging functionality
4. Add admin user management APIs
5. Create integration tests

## Support

See README.md for detailed documentation.
