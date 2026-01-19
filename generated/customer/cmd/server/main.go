package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/generated/customer/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/generated/customer/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer-service/v1"
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

	customerRepo := postgres.NewCustomerRepository(db)
	addressRepo := postgres.NewAddressRepository(db)
	cartRepo := postgres.NewCartRepository(db)
	favoriteRepo := postgres.NewFavoriteRepository(db)
	paymentMethodRepo := postgres.NewPaymentMethodRepository(db)
	reviewRepo := postgres.NewReviewRepository(db)

	getCustomerProfileUsecase := usecase.NewGetCustomerProfileUsecase(customerRepo)
	updateCustomerProfileUsecase := usecase.NewUpdateCustomerProfileUsecase(customerRepo)
	registerAddressUsecase := usecase.NewRegisterAddressUsecase(customerRepo, addressRepo)
	updateAddressUsecase := usecase.NewUpdateAddressUsecase(addressRepo)
	deleteAddressUsecase := usecase.NewDeleteAddressUsecase(addressRepo)
	addToCartUsecase := usecase.NewAddToCartUsecase(cartRepo)
	getCartUsecase := usecase.NewGetCartUsecase(cartRepo)
	updateCartItemQuantityUsecase := usecase.NewUpdateCartItemQuantityUsecase(cartRepo)
	removeFromCartUsecase := usecase.NewRemoveFromCartUsecase(cartRepo)
	mergeGuestCartUsecase := usecase.NewMergeGuestCartUsecase(cartRepo)
	addToFavoriteUsecase := usecase.NewAddToFavoriteUsecase(favoriteRepo)
	listFavoritesUsecase := usecase.NewListFavoritesUsecase(favoriteRepo)
	removeFromFavoriteUsecase := usecase.NewRemoveFromFavoriteUsecase(favoriteRepo)
	addPaymentMethodUsecase := usecase.NewAddPaymentMethodUsecase(paymentMethodRepo)
	deletePaymentMethodUsecase := usecase.NewDeletePaymentMethodUsecase(paymentMethodRepo)
	createReviewUsecase := usecase.NewCreateReviewUsecase(reviewRepo)
	updateReviewUsecase := usecase.NewUpdateReviewUsecase(reviewRepo)

	handler := grpchandler.NewCustomerServiceHandler(
		getCustomerProfileUsecase,
		updateCustomerProfileUsecase,
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

	log.Println("✅ All usecases and handlers initialized successfully")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCustomerServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("🚀 Customer Service gRPC server listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
