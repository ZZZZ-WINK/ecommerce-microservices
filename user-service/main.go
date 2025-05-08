package main

import (
	"common/middleware"
	"log"
	"net"
	"user-service/model"
	"user-service/service"

	pb "common/proto/gen/user"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:123mysql@tcp(127.0.0.1:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":50051")
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
