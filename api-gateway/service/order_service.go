package service

import (
	"context"
	"log"

	"api-gateway/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderService struct {
	conn   *grpc.ClientConn
	client proto.OrderServiceClient
}

func NewOrderService(address string) (*OrderService, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to order service: %v", err)
		return nil, err
	}

	client := proto.NewOrderServiceClient(conn)
	return &OrderService{
		conn:   conn,
		client: client,
	}, nil
}

func (s *OrderService) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	return s.client.CreateOrder(ctx, req)
}

func (s *OrderService) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	return s.client.GetOrder(ctx, req)
}

func (s *OrderService) ListOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	return s.client.ListOrders(ctx, req)
}
