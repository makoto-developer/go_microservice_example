package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/makoto-developer/go_microservice_example/generated/shop/proto/shop_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/shop/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/generated/shop/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Initialize repositories
	shopRepo := postgres.NewShopRepository(db)
	productRepo := postgres.NewProductRepository(db)
	shopCategoryRepo := postgres.NewShopCategoryRepository(db)

	// Initialize use cases
	shopRegistrationUsecase := usecase.NewShopRegistrationUsecase(
		shopRepo,
		shopCategoryRepo,
	)

	productManagementUsecase := usecase.NewProductManagementUsecase(
		productRepo,
		shopRepo,
	)

	// Initialize gRPC handler
	handler := grpchandler.NewShopServiceHandler(
		shopRegistrationUsecase,
		productManagementUsecase,
		shopRepo,
	)

	// Create gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShopServiceServer(grpcServer, handler)

	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	log.Printf("Shop Service is running on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
