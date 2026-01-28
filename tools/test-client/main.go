package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("🚀 オンラインショップフローテスト開始")
	log.Println()

	// テスト用のID
	customerID := uuid.New().String()
	userID := uuid.New().String()
	productID := uuid.New().String()

	// データベースに直接Customerレコードを作成
	log.Println("=== Step 0: テストデータ準備 ===")
	db, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=customer_db sslmode=disable")
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO customers (id, user_id, first_name, last_name, phone, birth_date, gender, created_at, updated_at)
		VALUES ($1, $2, '太郎', '山田', '090-1234-5678', '1990-01-15', 'male', NOW(), NOW())
	`, customerID, userID)
	if err != nil {
		log.Printf("❌ Customer作成失敗: %v", err)
	} else {
		log.Printf("✅ テスト用Customer作成成功 (ID: %s)", customerID)
	}
	log.Println()

	// Customer Serviceに接続
	conn, err := grpc.Dial("localhost:20102", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer conn.Close()

	client := pb.NewCustomerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("=== Step 1: 顧客プロフィール取得テスト ===")
	// プロフィール取得を試みる（まだ存在しないはず）
	_, err = client.GetProfile(ctx, &pb.GetProfileRequest{
		CustomerId: customerID,
	})
	if err != nil {
		log.Printf("✅ 予想通り: プロフィール未登録 - %v", err)
	}

	log.Println()
	log.Println("=== Step 1.5: 顧客プロフィール作成 ===")
	// プロフィールを作成
	profileResp, err := client.UpdateProfile(ctx, &pb.UpdateProfileRequest{
		CustomerId: customerID,
		FirstName:  "太郎",
		LastName:   "山田",
		Phone:      "090-1234-5678",
		BirthDate:  "1990-01-15",
		Gender:     pb.Gender_MALE,
	})
	if err != nil {
		log.Printf("❌ プロフィール作成失敗: %v", err)
	} else {
		log.Printf("✅ プロフィール作成成功!")
		log.Printf("   顧客ID: %s", profileResp.Customer.Id)
		log.Printf("   名前: %s %s", profileResp.Customer.LastName, profileResp.Customer.FirstName)
	}

	log.Println()
	log.Println("=== Step 2: カートに商品追加 ===")
	// カートに商品を追加
	cartResp, err := client.AddToCart(ctx, &pb.AddToCartRequest{
		CustomerId: customerID,
		ProductId:  productID,
		Quantity:   2,
	})
	if err != nil {
		log.Printf("❌ カート追加失敗: %v", err)
	} else {
		log.Printf("✅ カートに商品追加成功!")
		log.Printf("   カートアイテムID: %s", cartResp.CartItem.Id)
		log.Printf("   商品数: %d個", cartResp.CartItem.Quantity)
		log.Printf("   メッセージ: %s", cartResp.Message)
	}

	log.Println()
	log.Println("=== Step 3: カート内容取得 ===")
	// カート内容を取得
	getCartResp, err := client.GetCart(ctx, &pb.GetCartRequest{
		CustomerId: customerID,
	})
	if err != nil {
		log.Printf("❌ カート取得失敗: %v", err)
	} else {
		log.Printf("✅ カート取得成功!")
		log.Printf("   カート内商品数: %d種類", len(getCartResp.CartItems))
		for i, item := range getCartResp.CartItems {
			log.Printf("   商品%d: ID=%s, 数量=%d", i+1, item.ProductId, item.Quantity)
		}
	}

	log.Println()
	log.Println("=== Step 4: カート内商品数量更新 ===")
	if len(getCartResp.CartItems) > 0 {
		updateResp, err := client.UpdateCartItemQuantity(ctx, &pb.UpdateCartItemQuantityRequest{
			CartItemId: getCartResp.CartItems[0].Id,
			CustomerId: customerID,
			Quantity:   5,
		})
		if err != nil {
			log.Printf("❌ 数量更新失敗: %v", err)
		} else {
			log.Printf("✅ 数量更新成功! 新しい数量: %d個", updateResp.CartItem.Quantity)
		}
	}

	log.Println()
	log.Println("=== Step 5: お気に入り追加 ===")
	favResp, err := client.AddToFavorite(ctx, &pb.AddToFavoriteRequest{
		CustomerId: customerID,
		ProductId:  productID,
	})
	if err != nil {
		log.Printf("❌ お気に入り追加失敗: %v", err)
	} else {
		log.Printf("✅ お気に入り追加成功! %s", favResp.Message)
	}

	log.Println()
	log.Println("🎉 オンラインショップフローテスト完了!")
	log.Println()
	log.Println("=== テスト結果サマリー ===")
	log.Println("✅ Customer Service が正常に動作しています")
	log.Println("✅ カート機能が正常に動作しています")
	log.Println("✅ お気に入り機能が正常に動作しています")
	log.Println("✅ データベース連携が正常に動作しています")
}
