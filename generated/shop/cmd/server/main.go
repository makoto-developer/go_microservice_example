package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/generated/shop/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/generated/shop/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/shop_service/v1"
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

	shopRepo := postgres.NewShopRepository(db)
	productRepo := postgres.NewProductRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	salesRepo := postgres.NewSalesRepository(db)

	registerShopUsecase := usecase.NewRegisterShopUsecase(shopRepo)
	updateShopUsecase := usecase.NewUpdateShopUsecase(shopRepo)
	toggleShopPublishUsecase := usecase.NewToggleShopPublishUsecase(shopRepo)
	getShopUsecase := usecase.NewGetShopUsecase(shopRepo)

	registerProductUsecase := usecase.NewRegisterProductUsecase(shopRepo, productRepo)
	updateProductUsecase := usecase.NewUpdateProductUsecase(productRepo)
	deleteProductUsecase := usecase.NewDeleteProductUsecase(productRepo)
	toggleProductPublishUsecase := usecase.NewToggleProductPublishUsecase(productRepo)
	getProductUsecase := usecase.NewGetProductUsecase(productRepo)
	listProductsUsecase := usecase.NewListProductsUsecase(productRepo)
	uploadProductImageUsecase := usecase.NewUploadProductImageUsecase(productRepo)
	manageVariationUsecase := usecase.NewManageVariationUsecase(productRepo)

	listOrdersUsecase := usecase.NewListOrdersUsecase(orderRepo)
	updateOrderStatusUsecase := usecase.NewUpdateOrderStatusUsecase(orderRepo)

	getSalesReportUsecase := usecase.NewGetSalesReportUsecase(shopRepo, salesRepo)

	handler := grpchandler.NewShopServiceHandler(
		registerShopUsecase,
		updateShopUsecase,
		toggleShopPublishUsecase,
		getShopUsecase,
		registerProductUsecase,
		updateProductUsecase,
		deleteProductUsecase,
		toggleProductPublishUsecase,
		getProductUsecase,
		listProductsUsecase,
		uploadProductImageUsecase,
		manageVariationUsecase,
		listOrdersUsecase,
		updateOrderStatusUsecase,
		getSalesReportUsecase,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShopServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("🚀 Shop Service gRPC server listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
