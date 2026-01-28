package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	custompb "github.com/makoto-developer/go_microservice_example/proto/customer-service/v1"
	orderpb "github.com/makoto-developer/go_microservice_example/proto/order_service/v1"
	paymentpb "github.com/makoto-developer/go_microservice_example/proto/payment_service/v1"
	shoppb "github.com/makoto-developer/go_microservice_example/proto/shop_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("🎉 完全購入フローテスト開始")
	log.Println()

	ctx := context.Background()

	// テスト用のID
	shopOwnerID := uuid.New().String()
	customerID := uuid.New().String()
	userID := uuid.New().String()

	// Step 0: データベース準備
	log.Println("=== Step 0: テストデータ準備 ===")
	
	// Customer DB準備
	customerDB, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=customer_db sslmode=disable")
	if err != nil {
		log.Fatalf("Customer DB接続失敗: %v", err)
	}
	defer customerDB.Close()

	_, err = customerDB.Exec(`
		INSERT INTO customers (id, user_id, first_name, last_name, phone, birth_date, gender, created_at, updated_at)
		VALUES ($1, $2, '太郎', '山田', '090-1234-5678', '1990-01-15', 'male', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, customerID, userID)
	if err != nil {
		log.Printf("❌ Customer作成失敗: %v", err)
	} else {
		log.Printf("✅ テスト用Customer作成成功 (ID: %s)", customerID)
	}
	log.Println()

	// Step 1: Shop Serviceに接続して商品登録
	log.Println("=== Step 1: ショップと商品の登録 ===")
	
	shopConn, err := grpc.Dial("localhost:20101", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Shop Service接続失敗: %v", err)
	}
	defer shopConn.Close()

	shopClient := shoppb.NewShopServiceClient(shopConn)
	shopCtx, cancel1 := context.WithTimeout(ctx, 10*time.Second)
	defer cancel1()

	// ショップ登録
	shopResp, err := shopClient.RegisterShop(shopCtx, &shoppb.RegisterShopRequest{
		OwnerId:       shopOwnerID,
		Name:          "テストショップ",
		Description:   "テスト用のオンラインショップです",
		OwnerName:     "店主 花子",
		PhoneNumber:   "03-1234-5678",
		BusinessHours: "9:00-18:00",
		ReturnPolicy:  "30日以内返品可能",
	})
	if err != nil {
		log.Fatalf("❌ ショップ登録失敗: %v", err)
	}
	shopID := shopResp.ShopId
	log.Printf("✅ ショップ登録成功! Shop ID: %s", shopID)

	// 商品登録
	productResp, err := shopClient.RegisterProduct(shopCtx, &shoppb.RegisterProductRequest{
		ShopId:        shopID,
		Name:          "テスト商品A",
		Description:   "これはテスト用の商品です",
		Price:         "1000",
		Category:      "electronics",
		StockQuantity: 100,
	})
	if err != nil {
		log.Fatalf("❌ 商品登録失敗: %v", err)
	}
	productID := productResp.ProductId
	log.Printf("✅ 商品登録成功! Product ID: %s (価格: ¥1000, 在庫: 100個)", productID)
	log.Println()

	// Step 2: Customer Serviceに接続してカート追加
	log.Println("=== Step 2: カートに商品追加 ===")
	
	customerConn, err := grpc.Dial("localhost:20102", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Customer Service接続失敗: %v", err)
	}
	defer customerConn.Close()

	customerClient := custompb.NewCustomerServiceClient(customerConn)
	customerCtx, cancel2 := context.WithTimeout(ctx, 10*time.Second)
	defer cancel2()

	cartResp, err := customerClient.AddToCart(customerCtx, &custompb.AddToCartRequest{
		CustomerId: customerID,
		ProductId:  productID,
		Quantity:   3,
	})
	if err != nil {
		log.Fatalf("❌ カート追加失敗: %v", err)
	}
	log.Printf("✅ カート追加成功! 商品数: %d個", cartResp.CartItem.Quantity)
	log.Println()

	// Step 3: Order Serviceに接続して注文作成
	log.Println("=== Step 3: 注文作成 ===")
	
	orderConn, err := grpc.Dial("localhost:20104", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Order Service接続失敗: %v", err)
	}
	defer orderConn.Close()

	orderClient := orderpb.NewOrderServiceClient(orderConn)
	orderCtx, cancel3 := context.WithTimeout(ctx, 10*time.Second)
	defer cancel3()

	orderResp, err := orderClient.CreateOrder(orderCtx, &orderpb.CreateOrderRequest{
		CustomerId: customerID,
		CartItems: []*orderpb.CartItemInput{
			{
				ProductId: productID,
				Quantity:  3,
				UnitPrice: 1000,
			},
		},
	})
	if err != nil {
		log.Fatalf("❌ 注文作成失敗: %v", err)
	}
	orderID := orderResp.OrderId
	log.Printf("✅ 注文作成成功!")
	log.Printf("   注文ID: %s", orderID)
	log.Printf("   注文番号: %s", orderResp.OrderNumber)
	log.Println()

	// Step 4: Payment Serviceに接続して決済処理
	log.Println("=== Step 4: 決済処理 ===")
	
	paymentConn, err := grpc.Dial("localhost:20105", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Payment Service接続失敗: %v", err)
	}
	defer paymentConn.Close()

	paymentClient := paymentpb.NewPaymentServiceClient(paymentConn)
	paymentCtx, cancel4 := context.WithTimeout(ctx, 10*time.Second)
	defer cancel4()

	paymentResp, err := paymentClient.CreatePaymentIntent(paymentCtx, &paymentpb.CreatePaymentIntentRequest{
		OrderId:         orderID,
		Amount:          "3000",
		Currency:        "JPY",
		PaymentMethodId: "pm_card_visa",
		CustomerId:      customerID,
	})
	if err != nil {
		log.Fatalf("❌ 決済失敗: %v", err)
	}
	log.Printf("✅ 決済Intent作成完了!")
	log.Printf("   決済ID: %s", paymentResp.PaymentId)
	log.Printf("   メッセージ: %s", paymentResp.Message)
	log.Println()

	// 完了サマリー
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🎊 完全購入フロー テスト完了！")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println()
	log.Println("✅ ショップ登録 → 商品登録")
	log.Println("✅ カート追加")
	log.Println("✅ 注文作成")
	log.Println("✅ 決済処理")
	log.Println()
	log.Println("すべてのマイクロサービスが正常に連携しています！")
}
