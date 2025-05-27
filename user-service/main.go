package main

import (
	"common/middleware"
	"fmt"
	"log"
	"net"
	"os"
	"user-service/model"
	"user-service/service"

	pb "common/proto/gen/user"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 创建 gRPC 服务器
	port := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建带有 JWT 中间件的 gRPC 服务器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTMiddleware),
	)
	userService := service.NewUserService(db)
	pb.RegisterUserServiceServer(s, userService)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
