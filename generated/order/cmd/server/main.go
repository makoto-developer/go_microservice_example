package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/makoto-developer/go_microservice_example/proto/order_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/order/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/generated/order/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/usecase"
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

	log.Println("Successfully connected to database")

	orderRepo := postgres.NewOrderRepository(db)
	orderItemRepo := postgres.NewOrderItemRepository(db)
	orderMgmt := usecase.NewOrderManagementUsecase(orderRepo, orderItemRepo)
	handler := grpchandler.NewOrderServiceHandler(orderMgmt)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("Order Service is running on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
