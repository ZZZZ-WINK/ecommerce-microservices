package main

import (
	"context"
	"log"
	"time"

	pb "user-service/common/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接服务器: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := pb.NewUserServiceClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 测试注册功能
	registerResp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: "testuser",
		Password: "testpass",
		Email:    "test@example.com",
	})
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}
	log.Printf("注册成功: %v", registerResp)

	// 测试登录功能
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	})
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	log.Printf("登录成功: %v", loginResp)
}
