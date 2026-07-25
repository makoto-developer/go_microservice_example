// Package interceptor は order サービスの gRPC サーバに差し込む横断的な
// インターセプタ(認証・ロギング・パニック回復)を提供する。
//
// これらは grpc.NewServer(grpc.ChainUnaryInterceptor(...)) で登録され、
// 各 RPC ハンドラの「前段」で全リクエストに対して実行される。認可のような
// 横断的関心事はハンドラを直接呼ばないため、コールグラフには現れない。
package interceptor

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AccessClaims は auth サービスが HS256 で発行するアクセストークンの claims。
// サービス間はトークンの契約を共有する(本サンプルでは各サービスに複製している)。
type AccessClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"` // "customer" または "owner"
	jwt.RegisteredClaims
}

type ctxKey struct{}

// UserFromContext は認証済みユーザーの claims を context から取り出す。
// 認証インターセプタを通過した RPC のハンドラ内で利用できる。
func UserFromContext(ctx context.Context) (*AccessClaims, bool) {
	c, ok := ctx.Value(ctxKey{}).(*AccessClaims)
	return c, ok
}

// NewAuthUnaryInterceptor は `Authorization: Bearer <JWT>` を検証する認証
// インターセプタを返す。検証に成功すると claims を context に載せてハンドラへ渡す。
// トークンが無い/不正/期限切れの場合は codes.Unauthenticated を返す。
// gRPC reflection / health(メソッド名が "/grpc." で始まる)は認証不要で素通しする。
func NewAuthUnaryInterceptor(accessSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.HasPrefix(info.FullMethod, "/grpc.") {
			return handler(ctx, req)
		}

		raw, err := bearerToken(ctx)
		if err != nil {
			return nil, err
		}

		claims := &AccessClaims{}
		token, err := jwt.ParseWithClaims(
			raw,
			claims,
			func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, status.Error(codes.Unauthenticated, "unexpected signing method")
				}
				return []byte(accessSecret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithIssuer("auth-service"),
		)
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired access token")
		}

		return handler(context.WithValue(ctx, ctxKey{}, claims), req)
	}
}

// bearerToken は受信メタデータから Bearer トークン本体を取り出す。
func bearerToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing request metadata")
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization token")
	}
	parts := strings.SplitN(vals[0], " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", status.Error(codes.Unauthenticated, "authorization header must be 'Bearer <token>'")
	}
	return strings.TrimSpace(parts[1]), nil
}
