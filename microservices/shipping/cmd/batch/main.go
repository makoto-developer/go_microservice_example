// 配送メンテナンスバッチ。
// 配送業者 API から追跡情報を定期的に同期する(cron から起動される想定)。
// SHIPPING_BATCH_ONCE=1 で 1 回だけ実行して終了する。
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/makoto-developer/go_microservice_example/microservices/shipping/proto"
)

func main() {
	host := getEnv("SHIPPING_SERVICE_HOST", "localhost")
	port := getEnv("SHIPPING_SERVICE_PORT", "50057")
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to shipping service: %v", err)
	}
	defer conn.Close()
	stub := pb.NewShippingServiceClient(conn)

	interval := 10 * time.Minute
	log.Printf("🚚 Carrier tracking sync batch started (target: %s, interval: %s)", addr, interval)

	for {
		runOnce(stub)
		if os.Getenv("SHIPPING_BATCH_ONCE") == "1" {
			return
		}
		time.Sleep(interval)
	}
}

// runOnce は追跡情報の同期を 1 回実行する。
// 対象の出荷 ID は環境変数 SHIPPING_SYNC_TARGET(カンマ区切り)で指定する。
// 実運用では「輸送中の全出荷」をサービス側で列挙して同期する想定。
func runOnce(stub pb.ShippingServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	targets := os.Getenv("SHIPPING_SYNC_TARGET")
	if targets == "" {
		log.Printf("no sync targets configured (set SHIPPING_SYNC_TARGET)")
		return
	}
	for _, id := range splitComma(targets) {
		resp, err := stub.SyncCarrierTracking(ctx, &pb.SyncCarrierTrackingRequest{ShipmentId: id})
		if err != nil {
			log.Printf("sync failed for shipment %s: %v", id, err)
			continue
		}
		log.Printf("synced shipment %s: %s", id, resp.GetMessage())
	}
}

func splitComma(s string) []string {
	out := []string{}
	cur := ""
	for _, r := range s {
		if r == ',' {
			if cur != "" {
				out = append(out, cur)
			}
			cur = ""
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
