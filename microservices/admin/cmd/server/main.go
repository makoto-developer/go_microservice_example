package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/microservices/admin/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/microservices/admin/internal/handler/grpc"
	pb "github.com/makoto-developer/go_microservice_example/microservices/admin/proto"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Continuing with mock implementation...")
	} else {
		defer db.Close()
		if err := db.Ping(); err != nil {
			log.Printf("Warning: Failed to ping database: %v", err)
			log.Println("Continuing with mock implementation...")
		} else {
			log.Println("✅ Database connection established")
		}
	}

	handler := grpchandler.NewAdminServiceHandler()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("🚀 Admin Service gRPC server listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
