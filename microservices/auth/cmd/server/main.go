package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	
	"github.com/makoto-developer/go_microservice_example/generated/auth/config"
	grpchandler "github.com/makoto-developer/go_microservice_example/generated/auth/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/repository/postgres"
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/auth_service/v1"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	
	log.Println("Connected to database successfully")
	
	userRepo := postgres.NewUserRepository(db)
	jwtService := usecase.NewJWTService(cfg.JWTAccessSecret, cfg.JWTRefreshSecret)
	registrationUsecase := usecase.NewUserRegistrationUsecase(userRepo, jwtService)
	loginUsecase := usecase.NewUserLoginUsecase(userRepo, jwtService)
	authHandler := grpchandler.NewAuthServiceHandler(registrationUsecase, loginUsecase, jwtService)
	
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)
	
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ServerPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	log.Printf("Auth Service listening on port %s", cfg.ServerPort)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
