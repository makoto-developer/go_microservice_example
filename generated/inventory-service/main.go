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

	"github.com/makoto-developer/go_microservice_example/generated/inventory-service/handler"
	"github.com/makoto-developer/go_microservice_example/generated/inventory-service/infrastructure"
	"github.com/makoto-developer/go_microservice_example/generated/inventory-service/usecase"
	pb "github.com/makoto-developer/go_microservice_example/gen/inventory/v1"
)

func main() {
	// 環境変数からの設定読み込み
	grpcPort := getEnv("GRPC_PORT", "50052")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "admin")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "inventory_db")
	amqpURL := getEnv("AMQP_URL", "amqp://admin:password@localhost:5672/")

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

	// RabbitMQ接続
	eventPublisher, err := infrastructure.NewRabbitMQEventPublisher(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer eventPublisher.Close()
	log.Println("Successfully connected to RabbitMQ")

	// Repository層の初期化
	inventoryRepo := infrastructure.NewPostgresInventoryRepository(db)
	reservationRepo := infrastructure.NewPostgresReservationRepository(db)
	historyRepo := infrastructure.NewPostgresInventoryHistoryRepository(db)
	stockTakingRepo := infrastructure.NewPostgresStockTakingRepository(db)

	// UseCase層の初期化
	registerInventoryUC := usecase.NewRegisterInventory(inventoryRepo, historyRepo)
	updateInventoryUC := usecase.NewUpdateInventoryQuantity(inventoryRepo, historyRepo, eventPublisher)
	getInventoryUC := usecase.NewGetInventory(inventoryRepo)
	getInventoryByProductUC := usecase.NewGetInventoryByProduct(inventoryRepo)
	bulkGetInventoryUC := usecase.NewBulkGetInventory(inventoryRepo)
	reserveStockUC := usecase.NewReserveStock(inventoryRepo, reservationRepo, historyRepo, eventPublisher)
	bulkReserveStockUC := usecase.NewBulkReserveStock(inventoryRepo, reservationRepo, historyRepo, eventPublisher)
	releaseStockUC := usecase.NewReleaseStock(inventoryRepo, reservationRepo, historyRepo, eventPublisher)
	confirmStockUC := usecase.NewConfirmStock(inventoryRepo, reservationRepo, historyRepo, eventPublisher)
	getInventoryHistoryUC := usecase.NewGetInventoryHistory(historyRepo)
	recordStockTakingUC := usecase.NewRecordStockTaking(inventoryRepo, stockTakingRepo, historyRepo, eventPublisher)
	getStockTakingHistoryUC := usecase.NewGetStockTakingHistory(stockTakingRepo)

	// gRPC Handler初期化
	grpcHandler := handler.NewInventoryServiceHandler(
		registerInventoryUC,
		updateInventoryUC,
		getInventoryUC,
		getInventoryByProductUC,
		bulkGetInventoryUC,
		reserveStockUC,
		bulkReserveStockUC,
		releaseStockUC,
		confirmStockUC,
		getInventoryHistoryUC,
		recordStockTakingUC,
		getStockTakingHistoryUC,
	)

	// gRPCサーバー起動
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, grpcHandler)

	// Reflection有効化（開発環境のみ推奨）
	reflection.Register(s)

	// Graceful shutdown
	go func() {
		log.Printf("Inventory Service listening on :%s", grpcPort)
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
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
