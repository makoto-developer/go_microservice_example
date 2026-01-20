package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/makoto-developer/go_microservice_example/proto/customer_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/customer/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/generated/customer/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
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
