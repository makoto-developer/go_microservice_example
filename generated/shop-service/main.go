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

	"github.com/makoto-developer/go_microservice_example/generated/shop-service/handler"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/infrastructure"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/usecase"
	pb "github.com/makoto-developer/go_microservice_example/generated/shop-service/proto/shop_service/v1"
)

func main() {
	log.Println("Shop Service starting - BUILD VERSION: 2026-01-11-18:54")

	// 環境変数からの設定読み込み
	grpcPort := getEnv("GRPC_PORT", "20101")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "admin")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "shop_service_db")

	// Redis設定
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	// RabbitMQ設定
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://admin:password@localhost:5672/")

	// S3設定（画像アップロード用）
	s3Endpoint := getEnv("S3_ENDPOINT", "http://localhost:9000")
	s3AccessKey := getEnv("S3_ACCESS_KEY", "minioadmin")
	s3SecretKey := getEnv("S3_SECRET_KEY", "minioadmin")
	s3Bucket := getEnv("S3_BUCKET", "shop-images")

	log.Printf("Redis: %s:%s (password: %s)", redisHost, redisPort, maskSecret(redisPassword))
	log.Printf("RabbitMQ: %s", maskURL(rabbitmqURL))
	log.Printf("S3: %s (bucket: %s)", s3Endpoint, s3Bucket)

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

	// TODO: Redis接続（将来実装）
	_ = redisHost
	_ = redisPort

	// TODO: RabbitMQ接続（将来実装）
	// eventPublisher, err := infrastructure.NewRabbitMQEventPublisher(rabbitmqURL)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	// }
	// defer eventPublisher.Close()
	// log.Println("Successfully connected to RabbitMQ")

	// TODO: S3接続（将来実装）
	_ = s3AccessKey
	_ = s3SecretKey

	// Repository層の初期化
	shopRepo := infrastructure.NewPostgresShopRepository(db)
	shopCategoryRepo := infrastructure.NewPostgresShopCategoryRepository(db)
	productRepo := infrastructure.NewPostgresProductRepository(db)
	productImageRepo := infrastructure.NewPostgresProductImageRepository(db)
	productTagRepo := infrastructure.NewPostgresProductTagRepository(db)
	productVariationRepo := infrastructure.NewPostgresProductVariationRepository(db)
	orderRepo := infrastructure.NewPostgresOrderRepository(db)
	orderItemRepo := infrastructure.NewPostgresOrderItemRepository(db)
	salesReportRepo := infrastructure.NewPostgresSalesReportRepository(db)

	// UseCase層の初期化
	registerShopUC := usecase.NewRegisterShopUsecase(shopRepo, shopCategoryRepo)
	updateShopInfoUC := usecase.NewUpdateShopInfoUsecase(shopRepo)
	toggleShopPublishUC := usecase.NewToggleShopPublishUsecase(shopRepo)
	registerProductUC := usecase.NewRegisterProductUsecase(productRepo, productTagRepo)
	updateProductUC := usecase.NewUpdateProductUsecase(productRepo)
	deleteProductUC := usecase.NewDeleteProductUsecase(productRepo)
	toggleProductPublishUC := usecase.NewToggleProductPublishUsecase(productRepo)
	uploadProductImageUC := usecase.NewUploadProductImageUsecase(productImageRepo)
	getProductUC := usecase.NewGetProductUsecase(productRepo)
	listProductsUC := usecase.NewListProductsUsecase(productRepo)
	manageProductVariationUC := usecase.NewManageProductVariationUsecase(productVariationRepo)
	listOrdersUC := usecase.NewListOrdersUsecase(orderRepo)
	getOrderDetailUC := usecase.NewGetOrderDetailUsecase(orderRepo, orderItemRepo)
	updateOrderStatusUC := usecase.NewUpdateOrderStatusUsecase(orderRepo)
	getSalesReportUC := usecase.NewGetSalesReportUsecase(salesReportRepo)
	exportSalesDataUC := usecase.NewExportSalesDataUsecase(salesReportRepo)

	// gRPC Handler初期化
	grpcHandler := handler.NewShopServiceHandler(
		registerShopUC,
		updateShopInfoUC,
		toggleShopPublishUC,
		registerProductUC,
		updateProductUC,
		deleteProductUC,
		toggleProductPublishUC,
		uploadProductImageUC,
		getProductUC,
		listProductsUC,
		manageProductVariationUC,
		listOrdersUC,
		getOrderDetailUC,
		updateOrderStatusUC,
		getSalesReportUC,
		exportSalesDataUC,
	)

	// gRPCサーバー起動
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterShopServiceServer(s, grpcHandler)

	// Reflection有効化（開発環境のみ推奨）
	reflection.Register(s)

	// Graceful shutdown
	go func() {
		log.Printf("Shop Service listening on :%s", grpcPort)
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

// maskURL は接続URLをマスクする
func maskURL(url string) string {
	if url == "" {
		return "(empty)"
	}
	// amqp://user:password@host:port/ のpassword部分をマスク
	// 簡易的な実装
	if len(url) > 20 {
		return url[:10] + "****" + url[len(url)-10:]
	}
	return "****"
}
