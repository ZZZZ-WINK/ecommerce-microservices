package main

import (
	"api-gateway/config"
	"api-gateway/router"
	"api-gateway/service"
	"log"
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
		log.Fatal("Failed to connect to user service: ", err)
	}
	defer userSvc.Close()

	r := router.SetupRouter(userSvc)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
