package main

import (
	"common/middleware"
	"log"
	"net"
	"order-service/model"
	"order-service/service"

	pb "common/proto/gen/order"
	pbProduct "common/proto/gen/product"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123mysql@tcp(127.0.0.1:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 连接商品服务
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to product service: %v", err)
	}
	defer productConn.Close()
	productClient := pbProduct.NewProductServiceClient(productConn)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建带有 JWT 中间件的 gRPC 服务器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTMiddleware),
	)
	orderService := service.NewOrderService(db, productClient)
	pb.RegisterOrderServiceServer(s, orderService)

	log.Printf("Order service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
