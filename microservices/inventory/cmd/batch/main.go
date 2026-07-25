// 在庫メンテナンスバッチ。
// 期限切れの在庫引当を定期的に解放する(cron から起動される想定)。
// INVENTORY_BATCH_ONCE=1 で 1 回だけ実行して終了する。
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/makoto-developer/go_microservice_example/microservices/inventory/proto"
)

func main() {
	host := getEnv("INVENTORY_SERVICE_HOST", "localhost")
	port := getEnv("INVENTORY_SERVICE_PORT", "50054")
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer conn.Close()
	stub := pb.NewInventoryServiceClient(conn)

	interval := 5 * time.Minute
	log.Printf("🧹 Inventory maintenance batch started (target: %s, interval: %s)", addr, interval)

	for {
		runOnce(stub)
		if os.Getenv("INVENTORY_BATCH_ONCE") == "1" {
			return
		}
		time.Sleep(interval)
	}
}

// runOnce は期限切れ引当の解放を 1 回実行する。
func runOnce(stub pb.InventoryServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := stub.ReleaseExpiredReservations(ctx, &pb.ReleaseExpiredReservationsRequest{})
	if err != nil {
		log.Printf("release expired reservations failed: %v", err)
		return
	}
	log.Printf("released expired reservations: %s (count=%d)", resp.GetMessage(), resp.GetReleasedCount())
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
