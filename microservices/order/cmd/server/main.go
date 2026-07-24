package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/microservices/order/config"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/client"
	grpchandler "github.com/makoto-developer/go_microservice_example/microservices/order/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/order/proto"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("✅ Successfully connected to Order database at %s:%s", cfg.Database.Host, cfg.Database.Port)

	orderRepo := postgres.NewOrderRepository(db)
	orderItemRepo := postgres.NewOrderItemRepository(db)

	paymentClient, err := client.NewPaymentClient(cfg.Payment.Address())
	if err != nil {
		log.Fatalf("Failed to create payment client: %v", err)
	}
	defer paymentClient.Close()
	log.Printf("💳 Payment Service: %s", cfg.Payment.Address())

	orderMgmt := usecase.NewOrderManagementUsecase(orderRepo, orderItemRepo, paymentClient)
	handler := grpchandler.NewOrderServiceHandler(orderMgmt)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("✅ Order Service is running on port %s", cfg.Server.Port)
	log.Printf("🎯 Database per Service architecture is active!")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
