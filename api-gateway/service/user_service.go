package service

import (
	"context"
	"log"

	"api-gateway/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserService struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
}

func NewUserService(address string) (*UserService, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
		return nil, err
	}

	client := proto.NewUserServiceClient(conn)
	return &UserService{
		conn:   conn,
		client: client,
	}, nil
}

func (s *UserService) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *UserService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return s.client.Register(ctx, req)
}

func (s *UserService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return s.client.Login(ctx, req)
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	return s.client.GetUserInfo(ctx, req)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	return s.client.UpdateUser(ctx, req)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return s.client.DeleteUser(ctx, req)
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	return s.client.ListUsers(ctx, req)
}
