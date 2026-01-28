#!/bin/bash

# Admin Service Database Setup Script

echo "Setting up Admin Service database..."

# Database connection details
DB_HOST="localhost"
DB_PORT="22021"
DB_NAME="admin_service"
DB_USER="postgres"
DB_PASSWORD="postgres_password"

# Set PGPASSWORD for psql
export PGPASSWORD=$DB_PASSWORD

# Execute schema
echo "Executing schema.sql..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f sql/schema.sql

if [ $? -eq 0 ]; then
    echo "✅ Admin Service database setup completed successfully!"
else
    echo "❌ Failed to setup Admin Service database"
    exit 1
fi
