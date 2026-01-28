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
	// Database configuration - use same format as Auth Service
	databaseURL := getEnv("ORDER_DATABASE_URL", "postgresql://postgres:postgres_password@localhost:22014/order_service?sslmode=disable")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("✅ Successfully connected to Order database")

	// Start gRPC server
	serverPort := getEnv("ORDER_SERVICE_PORT", "22104")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	log.Printf("✅ Order Service is running on port %s", serverPort)
	log.Printf("🎯 Database per Service architecture is active!")
	log.Printf("   - Order Service has dedicated PostgreSQL instance on port 22014")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
