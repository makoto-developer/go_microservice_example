# Search Service

## Overview
Search Service provides full-text search functionality across products, shops, and other searchable entities in the microservices architecture.

## Features
- Full-text search indexing
- Search query logging and analytics
- Multi-entity search support (products, shops, etc.)
- Search result ranking
- Search history tracking

## Database Schema

### Tables
1. **search_indexes** - Search index data
   - `id`: UUID primary key
   - `entity_type`: Type of entity (product, shop, etc.)
   - `entity_id`: Reference to the indexed entity
   - `searchable_text`: Full-text searchable content
   - `metadata`: Additional JSONB metadata
   - `created_at`, `updated_at`: Timestamps

2. **search_logs** - Search query logs
   - `id`: UUID primary key
   - `user_id`: User who performed the search
   - `query_text`: Search query string
   - `filters`: Applied filters (JSONB)
   - `result_count`: Number of results returned
   - `search_time_ms`: Query execution time
   - `created_at`: Timestamp

## Configuration

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string (default: `postgresql://postgres:postgres_password@localhost:22020/search_service?sslmode=disable`)
- `SERVICE_PORT`: gRPC server port (default: `22110`)

### Database
- **PostgreSQL Port**: 22020
- **Database Name**: search_service
- **Password**: postgres_password

## Running the Service

### Prerequisites
1. PostgreSQL must be running on port 22020
2. Go 1.23 or higher

### Start the Service
```bash
# Install dependencies
go mod download

# Run the service
go run main.go
```

### Using Docker Compose
```bash
# From infrastructure/docker directory
docker-compose up postgres_search -d

# Wait for database to be ready
docker-compose ps postgres_search

# Run the service
cd simple-servers/search
go run main.go
```

## Build

```bash
# Build binary
go build -o search-service main.go

# Run binary
./search-service
```

## Testing Database Connection

```bash
# Test PostgreSQL connection
psql -h localhost -p 22020 -U postgres -d search_service

# Verify tables
\dt

# Check search_indexes table
SELECT * FROM search_indexes LIMIT 10;

# Check search_logs table
SELECT * FROM search_logs LIMIT 10;
```

## Architecture Notes

This service follows the **Database per Service** pattern:
- Dedicated PostgreSQL instance on port 22020
- Independent schema management
- Isolated data storage
- Service-specific indexes and optimizations

## Integration Points

### Search Indexing
- Products from Shop Service
- Shop information from Shop Service
- Reviews from Review Service
- Customer data (when permitted)

### Search Queries
- Full-text search across indexed content
- Filter by entity type
- Sort by relevance or other criteria
- Pagination support

## Performance Considerations

1. **Indexing Strategy**
   - Indexes on entity_type and entity_id for fast lookups
   - Consider PostgreSQL full-text search capabilities
   - Future: Integrate with Elasticsearch for advanced search

2. **Query Optimization**
   - Log slow queries for analysis
   - Cache frequently searched terms
   - Implement query result caching

3. **Monitoring**
   - Track search query performance
   - Monitor index update frequency
   - Analyze popular search terms

## Future Enhancements
- Elasticsearch integration for advanced full-text search
- Search suggestions and autocomplete
- Personalized search results
- Search analytics dashboard
- Multi-language search support
