package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/makoto-developer/go_microservice_example/microservices/shipping/config"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/client"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/domain"
	grpchandler "github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/handler/grpc"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/shipping/proto"
)

// mockShipmentRepository is a simple in-memory repository for testing
type mockShipmentRepository struct {
	shipments map[uuid.UUID]*domain.Shipment
}

func newMockShipmentRepository() repository.ShipmentRepository {
	return &mockShipmentRepository{
		shipments: make(map[uuid.UUID]*domain.Shipment),
	}
}

func (r *mockShipmentRepository) Create(ctx context.Context, shipment *domain.Shipment) error {
	r.shipments[shipment.ID] = shipment
	return nil
}

func (r *mockShipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Shipment, error) {
	if shipment, exists := r.shipments[id]; exists {
		return shipment, nil
	}
	return nil, nil
}

func (r *mockShipmentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Shipment, error) {
	for _, shipment := range r.shipments {
		if shipment.OrderID == orderID {
			return shipment, nil
		}
	}
	return nil, nil
}

func (r *mockShipmentRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Shipment, error) {
	var result []*domain.Shipment
	for _, shipment := range r.shipments {
		if shipment.CustomerID == customerID {
			result = append(result, shipment)
		}
	}
	return result, nil
}

func (r *mockShipmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ShipmentStatus) error {
	if shipment, exists := r.shipments[id]; exists {
		shipment.Status = status
		return nil
	}
	return fmt.Errorf("shipment not found")
}

func (r *mockShipmentRepository) UpdateTracking(ctx context.Context, id uuid.UUID, trackingNumber string) error {
	if shipment, exists := r.shipments[id]; exists {
		shipment.TrackingNumber = trackingNumber
		return nil
	}
	return fmt.Errorf("shipment not found")
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Database connection
	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Continuing with mock repository...")
	} else {
		defer db.Close()
		if err := db.Ping(); err != nil {
			log.Printf("Warning: Failed to ping database: %v", err)
			log.Println("Continuing with mock repository...")
		} else {
			log.Println("✅ Database connection established")
		}
	}

	// Use mock repository for now
	shipmentRepo := newMockShipmentRepository()

	// Initialize usecases
	createShipmentUsecase := usecase.NewCreateShipmentUsecase(shipmentRepo)

	// 決済サービスへの接続(配達完了時の代引き集金確定に使う)
	paymentAddr := fmt.Sprintf("%s:%s", cfg.Payment.Host, cfg.Payment.Port)
	paymentClient, err := client.NewPaymentClient(paymentAddr)
	if err != nil {
		log.Printf("Warning: payment client unavailable (%v). COD confirmation disabled.", err)
		paymentClient = nil
	} else {
		defer paymentClient.Close()
		log.Printf("✅ Payment client configured for %s", paymentAddr)
	}

	// Initialize handler
	handler := grpchandler.NewShippingServiceHandler(
		createShipmentUsecase,
		shipmentRepo,
		paymentClient,
	)

	// 配達完了メール(通知サービス)
	notificationAddr := fmt.Sprintf("%s:%s", cfg.Notification.Host, cfg.Notification.Port)
	if notificationClient, err := client.NewNotificationClient(notificationAddr); err != nil {
		log.Printf("Warning: notification client unavailable (%v). Delivery emails disabled.", err)
	} else {
		defer notificationClient.Close()
		handler.WithNotification(notificationClient)
		log.Printf("✉️ Notification Service: %s", notificationAddr)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShippingServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("🚀 Shipping Service gRPC server listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
