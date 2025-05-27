package service

import (
	"context"
	"log"

	"api-gateway/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductService struct {
	conn   *grpc.ClientConn
	client proto.ProductServiceClient
}

func NewProductService(address string) (*ProductService, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to product service: %v", err)
		return nil, err
	}

	client := proto.NewProductServiceClient(conn)
	return &ProductService{
		conn:   conn,
		client: client,
	}, nil
}

func (s *ProductService) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *ProductService) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	return s.client.GetProduct(ctx, req)
}

func (s *ProductService) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	return s.client.ListProducts(ctx, req)
}

func (s *ProductService) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	return s.client.CreateProduct(ctx, req)
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	return s.client.UpdateProduct(ctx, req)
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	return s.client.DeleteProduct(ctx, req)
}
