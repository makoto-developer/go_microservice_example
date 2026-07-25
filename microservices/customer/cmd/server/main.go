package main

import (
	"database/sql"
	"fmt"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/client"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/microservices/customer/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/microservices/customer/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/customer/proto"
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

	// Initialize repositories
	customerRepo := postgres.NewCustomerRepository(db)
	addressRepo := postgres.NewAddressRepository(db)
	cartRepo := postgres.NewCartRepository(db)
	favoriteRepo := postgres.NewFavoriteRepository(db)
	paymentMethodRepo := postgres.NewPaymentMethodRepository(db)
	reviewRepo := postgres.NewReviewRepository(db)

	// Initialize profile usecases
	getProfileUsecase := usecase.NewGetCustomerProfileUsecase(customerRepo)
	updateProfileUsecase := usecase.NewUpdateCustomerProfileUsecase(customerRepo)

	// Initialize address usecases
	registerAddressUsecase := usecase.NewRegisterAddressUsecase(customerRepo, addressRepo)
	updateAddressUsecase := usecase.NewUpdateAddressUsecase(addressRepo)
	deleteAddressUsecase := usecase.NewDeleteAddressUsecase(addressRepo)

	// Initialize cart usecases
	addToCartUsecase := usecase.NewAddToCartUsecase(cartRepo)
	getCartUsecase := usecase.NewGetCartUsecase(cartRepo)
	updateCartItemQuantityUsecase := usecase.NewUpdateCartItemQuantityUsecase(cartRepo)
	removeFromCartUsecase := usecase.NewRemoveFromCartUsecase(cartRepo)
	mergeGuestCartUsecase := usecase.NewMergeGuestCartUsecase(cartRepo)

	// Initialize favorite usecases
	addToFavoriteUsecase := usecase.NewAddToFavoriteUsecase(favoriteRepo)
	listFavoritesUsecase := usecase.NewListFavoritesUsecase(favoriteRepo)
	removeFromFavoriteUsecase := usecase.NewRemoveFromFavoriteUsecase(favoriteRepo)

	// Initialize payment method usecases
	addPaymentMethodUsecase := usecase.NewAddPaymentMethodUsecase(paymentMethodRepo)
	deletePaymentMethodUsecase := usecase.NewDeletePaymentMethodUsecase(paymentMethodRepo)

	// Initialize review usecases
	createReviewUsecase := usecase.NewCreateReviewUsecase(reviewRepo)
	updateReviewUsecase := usecase.NewUpdateReviewUsecase(reviewRepo)

	// Initialize handler with all usecases
	handler := grpchandler.NewCustomerServiceHandler(
		getProfileUsecase,
		updateProfileUsecase,
		registerAddressUsecase,
		updateAddressUsecase,
		deleteAddressUsecase,
		addToCartUsecase,
		getCartUsecase,
		updateCartItemQuantityUsecase,
		removeFromCartUsecase,
		mergeGuestCartUsecase,
		addToFavoriteUsecase,
		listFavoritesUsecase,
		removeFromFavoriteUsecase,
		addPaymentMethodUsecase,
		deletePaymentMethodUsecase,
		createReviewUsecase,
		updateReviewUsecase,
	)

	// 注文履歴系(order サービス委譲)とレビュー参照の配線
	orderAddr := fmt.Sprintf("%s:%s", getEnvOr("ORDER_SERVICE_HOST", "localhost"), getEnvOr("ORDER_SERVICE_PORT", "50055"))
	if orderClient, err := client.NewOrderClient(orderAddr); err != nil {
		log.Printf("Warning: order client unavailable (%v). Order history disabled.", err)
	} else {
		defer orderClient.Close()
		handler.WithOrderClient(orderClient)
		log.Printf("🧾 Order Service: %s", orderAddr)
	}
	handler.WithReviewRepo(reviewRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCustomerServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("Customer Service is running on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getEnvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
