package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// NewLoggingUnaryInterceptor は各 RPC のメソッド名・所要時間・結果コードを記録する。
// 全 RPC の前段で実行される横断ロギング。
func NewLoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		log.Printf("grpc method=%s code=%s duration=%s", info.FullMethod, status.Code(err), time.Since(start))
		return resp, err
	}
}
