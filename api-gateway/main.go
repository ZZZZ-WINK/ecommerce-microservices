package main

import (
	"api-gateway/config"
	"api-gateway/router"
	"api-gateway/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title E-commerce API Gateway
// @version 1.0
// @description API Gateway for E-commerce Microservices
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 初始化 gRPC 用户服务客户端
	userSvc, err := service.NewUserService(config.DefaultConfig.Services.UserService)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userSvc.Close()

	// 初始化 gRPC 商品服务客户端
	productSvc, err := service.NewProductService(config.DefaultConfig.Services.ProductService)
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productSvc.Close()

	// TODO: 初始化 gRPC 订单服务客户端
	// orderSvc, err := service.NewOrderService(config.DefaultConfig.Services.OrderService)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to order service: %v", err)
	// }
	// defer orderSvc.Close()

	r := router.SetupRouter(userSvc, productSvc)

	// 启动 HTTP 服务器
	go func() {
		if err := r.Run(":" + config.DefaultConfig.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	// 在此处执行清理操作，例如关闭数据库连接等（如果不是由服务客户端处理）
}
