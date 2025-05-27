package service

import (
	pb "common/proto/gen/product"
	"context"
	"encoding/json"
	"fmt"
	"product-service/model"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		MainImage:   req.MainImage,
	}
	if err := s.db.Create(&product).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to create product")
	}
	return &pb.CreateProductResponse{
		ProductId: int64(product.ID),
		Message:   "product created successfully",
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	// 1. 先查 Redis
	cacheKey := fmt.Sprintf("product:detail:%d", req.ProductId)
	val, err := RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && val != "" {
		var cachedProduct pb.Product
		if json.Unmarshal([]byte(val), &cachedProduct) == nil {
			return &pb.GetProductResponse{Product: &cachedProduct}, nil
		}
	}

	// 2. 查数据库
	var product model.Product
	if err := s.db.First(&product, req.ProductId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to query product")
	}

	// 3. 写入 Redis
	pbProduct := &pb.Product{
		Id:          int64(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
		MainImage:   product.MainImage,
	}
	bytes, _ := json.Marshal(pbProduct)
	RedisClient.Set(ctx, cacheKey, bytes, 5*time.Minute)

	return &pb.GetProductResponse{Product: pbProduct}, nil
}

func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	// 只缓存无关键词、第一页的商品列表
	cacheKey := "product:list:page:1:size:10"
	if strings.TrimSpace(req.Keyword) == "" && (req.Page == 1 || req.Page == 0) && (req.PageSize == 10 || req.PageSize == 0) {
		val, err := RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			var cachedResp pb.ListProductsResponse
			if json.Unmarshal([]byte(val), &cachedResp) == nil {
				return &cachedResp, nil
			}
		}
	}

	var products []model.Product
	var total int64
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	dbQuery := s.db.Model(&model.Product{})
	if strings.TrimSpace(req.Keyword) != "" {
		kw := "%" + strings.TrimSpace(req.Keyword) + "%"
		dbQuery = dbQuery.Where("name LIKE ?", kw)
	}
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to count products")
	}
	if err := dbQuery.Limit(pageSize).Offset(offset).Order("id desc").Find(&products).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to list products")
	}
	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          int64(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       int32(p.Stock),
			MainImage:   p.MainImage,
		})
	}
	resp := &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(total),
	}
	// 写入缓存
	if strings.TrimSpace(req.Keyword) == "" && page == 1 && pageSize == 10 {
		bytes, _ := json.Marshal(resp)
		RedisClient.Set(ctx, cacheKey, bytes, 2*time.Minute)
	}
	return resp, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	var product model.Product
	if err := s.db.First(&product, req.ProductId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.UpdateProductResponse{Success: false, Message: "product not found"}, nil
		}
		return &pb.UpdateProductResponse{Success: false, Message: "failed to query product"}, nil
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Stock != 0 {
		product.Stock = int(req.Stock)
	}
	if req.MainImage != "" {
		product.MainImage = req.MainImage
	}
	if err := s.db.Save(&product).Error; err != nil {
		return &pb.UpdateProductResponse{Success: false, Message: "failed to update product"}, nil
	}
	return &pb.UpdateProductResponse{Success: true, Message: "product updated successfully"}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if err := s.db.Delete(&model.Product{}, req.ProductId).Error; err != nil {
		return &pb.DeleteProductResponse{Success: false, Message: "failed to delete product"}, nil
	}
	return &pb.DeleteProductResponse{Success: true, Message: "product deleted successfully"}, nil
}
