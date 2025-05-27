package middleware

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTSecret 密钥，从环境变量中读取
var JWTSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key"))

// getEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Claims 定义JWT的声明结构
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 不需要鉴权的接口白名单
var noAuthMethods = map[string]bool{
	"/proto.UserService/Register":        true,
	"/proto.UserService/Login":           true,
	"/proto.ProductService/ListProducts": true,
	"/proto.ProductService/GetProduct":   true,
}

// GenerateToken 生成JWT token
func GenerateToken(userID int64, username string) (string, error) {
	// 设置token的过期时间
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建Claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
			Subject:   username,
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// JWTMiddleware 是一个 gRPC 中间件，用于验证JWT token
func JWTMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 检查是否在白名单中
	if noAuthMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// 从上下文中获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// 从元数据中获取token
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	// 检查token格式
	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, status.Error(codes.Unauthenticated, "authorization token format is invalid")
	}

	tokenString := parts[1]

	// 解析token
	claims, err := ParseToken(tokenString)
	if err != nil {
		if err == ErrExpiredToken {
			return nil, status.Error(codes.Unauthenticated, "token has expired")
		}
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// 将用户信息存储到上下文中
	newCtx := context.WithValue(ctx, "user", claims)
	return handler(newCtx, req)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// 验证token
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// 获取Claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
