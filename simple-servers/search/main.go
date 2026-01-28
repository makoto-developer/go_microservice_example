package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Database configuration
	databaseURL := getEnv("DATABASE_URL", "postgresql://postgres:postgres_password@localhost:22020/search_service?sslmode=disable")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("✅ Successfully connected to Search database")

	// Initialize database schema
	if err := initializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	log.Printf("✅ Database schema initialized")

	// Start gRPC server
	serverPort := getEnv("SERVICE_PORT", "22110")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	log.Printf("✅ Search Service is running on port %s", serverPort)
	log.Printf("🎯 Database per Service architecture is active!")
	log.Printf("   - Search Service has dedicated PostgreSQL instance on port 22020")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func initializeSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS search_indexes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		entity_type VARCHAR(50) NOT NULL,
		entity_id UUID NOT NULL,
		searchable_text TEXT NOT NULL,
		metadata JSONB,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS search_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID,
		query_text VARCHAR(500) NOT NULL,
		filters JSONB,
		result_count INTEGER NOT NULL,
		search_time_ms INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_search_indexes_entity_type ON search_indexes(entity_type);
	CREATE INDEX IF NOT EXISTS idx_search_indexes_entity_id ON search_indexes(entity_id);
	CREATE INDEX IF NOT EXISTS idx_search_logs_user_id ON search_logs(user_id);
	CREATE INDEX IF NOT EXISTS idx_search_logs_created_at ON search_logs(created_at);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
