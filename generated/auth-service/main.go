package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/generated/auth-service/handler"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/infrastructure"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/usecase"
	pb "github.com/makoto-developer/go_microservice_example/generated/auth-service/proto/auth_service/v1"
	"github.com/makoto-developer/go_microservice_example/manual/auth"
)

func main() {
	// 環境変数からの設定読み込み
	grpcPort := getEnv("GRPC_PORT", "20100")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "admin")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "auth_service_db")

	// JWT設定（将来的に使用）
	jwtSecret := getEnv("JWT_SECRET", "default-secret-key-change-in-production")
	jwtAccessExpiration := getEnv("JWT_ACCESS_EXPIRATION", "15m")
	jwtRefreshExpiration := getEnv("JWT_REFRESH_EXPIRATION", "7d")

	// Redis設定（オプション、将来的にトークンキャッシュ用）
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	log.Printf("JWT Secret: %s (length: %d)", maskSecret(jwtSecret), len(jwtSecret))
	log.Printf("JWT Access Expiration: %s", jwtAccessExpiration)
	log.Printf("JWT Refresh Expiration: %s", jwtRefreshExpiration)
	log.Printf("Redis: %s:%s (password: %s)", redisHost, redisPort, maskSecret(redisPassword))

	// PostgreSQL接続
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	// Repository層の初期化
	userRepo := infrastructure.NewPostgresUserRepository(db)
	refreshTokenRepo := infrastructure.NewPostgresRefreshTokenRepository(db)

	// EmailSender の初期化
	emailSender := auth.NewSMTPEmailSender()

	// UseCase層の初期化
	userRegistrationUC := usecase.NewUserRegistrationUsecase(userRepo, refreshTokenRepo, jwtSecret, emailSender)
	emailVerificationUC := usecase.NewEmailVerificationUsecase(userRepo)
	userLoginUC := usecase.NewUserLoginUsecase(userRepo, refreshTokenRepo, jwtSecret)
	userLogoutUC := usecase.NewUserLogoutUsecase(refreshTokenRepo)
	tokenRefreshUC := usecase.NewTokenRefreshUsecase(refreshTokenRepo, userRepo, jwtSecret)
	tokenVerificationUC := usecase.NewTokenVerificationUsecase(userRepo, jwtSecret)
	passwordResetRequestUC := usecase.NewPasswordResetRequestUsecase(userRepo, emailSender)
	passwordResetUC := usecase.NewPasswordResetUsecase(userRepo)
	passwordChangeUC := usecase.NewPasswordChangeUsecase(userRepo)

	// gRPC Handler初期化
	grpcHandler := handler.NewAuthServiceHandler(
		userRegistrationUC,
		emailVerificationUC,
		userLoginUC,
		userLogoutUC,
		tokenRefreshUC,
		tokenVerificationUC,
		passwordResetRequestUC,
		passwordResetUC,
		passwordChangeUC,
	)

	// gRPCサーバー起動
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, grpcHandler)

	// Reflection有効化（開発環境のみ推奨）
	reflection.Register(s)

	// Graceful shutdown
	go func() {
		log.Printf("Auth Service listening on :%s", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server stopped")

	// Repositoryは未使用だが、将来の実装のため保持
	_ = userRepo
	_ = refreshTokenRepo
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// maskSecret は秘密情報を部分的にマスクする
func maskSecret(secret string) string {
	if secret == "" {
		return "(empty)"
	}
	if len(secret) <= 4 {
		return "****"
	}
	return secret[:2] + "****" + secret[len(secret)-2:]
}
