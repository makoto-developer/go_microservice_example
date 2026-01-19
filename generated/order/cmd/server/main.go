package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/generated/order/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/generated/order/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/order_service/v1"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✅ Database connection established")

	orderRepo := postgres.NewOrderRepository(db)

	createOrderUsecase := usecase.NewCreateOrderUsecase(orderRepo)
	cancelOrderUsecase := usecase.NewCancelOrderUsecase(orderRepo)

	handler := grpchandler.NewOrderServiceHandler(
		createOrderUsecase,
		cancelOrderUsecase,
		orderRepo,
	)

	log.Println("✅ All usecases and handlers initialized successfully")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("🚀 Order Service gRPC server listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
