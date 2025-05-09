package main

import (
	"common/middleware"
	"log"
	"net"
	"product-service/model"
	"product-service/service"

	pb "common/proto/gen/product"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "zli:123456@tcp(192.168.94.242:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 初始化 Redis
	service.InitRedis()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTMiddleware),
	)
	productService := service.NewProductService(db)
	pb.RegisterProductServiceServer(s, productService)

	log.Printf("Product service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
