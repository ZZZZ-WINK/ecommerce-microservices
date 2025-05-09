package service

import (
	pb "common/proto/gen/order"
	pbProduct "common/proto/gen/product"
	"context"
	"order-service/model"
	"time"

	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	db            *gorm.DB
	productClient pbProduct.ProductServiceClient
}

func NewOrderService(db *gorm.DB, productClient pbProduct.ProductServiceClient) *OrderService {
	return &OrderService{db: db, productClient: productClient}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// 1. 对每个商品加锁
	for _, item := range req.Items {
		lockKey := fmt.Sprintf("lock:product:%d", item.ProductId)
		ok, err := RedisClient.SetNX(Ctx, lockKey, "locked", 5*time.Second).Result()
		if err != nil || !ok {
			return nil, status.Error(codes.Aborted, "order is too frequent, please try again later")
		}
		defer RedisClient.Del(Ctx, lockKey) // 下单结束后自动释放锁
	}

	// 1. 检查库存，调用商品服务
	for _, item := range req.Items {
		resp, err := s.productClient.GetProduct(ctx, &pbProduct.GetProductRequest{ProductId: item.ProductId})
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to get product info: %v", err)
		}
		if resp.Product.Stock < item.Quantity {
			return nil, status.Errorf(codes.FailedPrecondition, "product %s stock not enough", resp.Product.Name)
		}
	}
	// 2. 计算总价
	total := 0.0
	for _, item := range req.Items {
		total += item.Price * float64(item.Quantity)
	}
	// 3. 创建订单和订单项
	order := model.Order{
		UserID:     req.UserId,
		TotalPrice: total,
		Status:     int(pb.OrderStatus_PENDING),
	}
	for _, item := range req.Items {
		order.Items = append(order.Items, model.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			Price:       item.Price,
		})
	}
	if err := s.db.Create(&order).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to create order")
	}

	// 发送 Kafka 消息
	msg := fmt.Sprintf("order_created|order_id=%d|user_id=%d", order.ID, order.UserID)
	err := SendOrderMessage(msg)
	if err != nil {
		fmt.Println("Failed to send Kafka message:", err)
	}
	return &pb.CreateOrderResponse{
		OrderId: int64(order.ID),
		Message: "order created successfully",
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	var order model.Order
	if err := s.db.Preload("Items").First(&order, req.OrderId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, "failed to query order")
	}
	return &pb.GetOrderResponse{
		Order: convertOrderModelToPB(&order),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	var orders []model.Order
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
	dbQuery := s.db.Model(&model.Order{}).Where("user_id = ?", req.UserId)
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to count orders")
	}
	if err := dbQuery.Preload("Items").Limit(pageSize).Offset(offset).Order("id desc").Find(&orders).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to list orders")
	}
	var pbOrders []*pb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, convertOrderModelToPB(&o))
	}
	return &pb.ListOrdersResponse{
		Orders: pbOrders,
		Total:  int32(total),
	}, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	var order model.Order
	if err := s.db.First(&order, req.OrderId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.UpdateOrderStatusResponse{Success: false, Message: "order not found"}, nil
		}
		return &pb.UpdateOrderStatusResponse{Success: false, Message: "failed to query order"}, nil
	}
	order.Status = int(req.Status)
	if err := s.db.Save(&order).Error; err != nil {
		return &pb.UpdateOrderStatusResponse{Success: false, Message: "failed to update status"}, nil
	}
	return &pb.UpdateOrderStatusResponse{Success: true, Message: "order status updated"}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	// 先查找订单是否存在
	var order model.Order
	if err := s.db.First(&order, req.OrderId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.DeleteOrderResponse{Success: false, Message: "order not found"}, nil
		}
		return &pb.DeleteOrderResponse{Success: false, Message: "failed to query order"}, nil
	}
	// 先删除订单项，再删订单
	if err := s.db.Where("order_id = ?", req.OrderId).Delete(&model.OrderItem{}).Error; err != nil {
		return &pb.DeleteOrderResponse{Success: false, Message: "failed to delete order items"}, nil
	}
	if err := s.db.Delete(&model.Order{}, req.OrderId).Error; err != nil {
		return &pb.DeleteOrderResponse{Success: false, Message: "failed to delete order"}, nil
	}
	return &pb.DeleteOrderResponse{Success: true, Message: "order deleted successfully"}, nil
}

// 工具函数：模型转pb
func convertOrderModelToPB(order *model.Order) *pb.Order {
	var items []*pb.OrderItem
	for _, it := range order.Items {
		items = append(items, &pb.OrderItem{
			ProductId:   it.ProductID,
			ProductName: it.ProductName,
			Quantity:    int32(it.Quantity),
			Price:       it.Price,
		})
	}
	return &pb.Order{
		Id:         int64(order.ID),
		UserId:     order.UserID,
		Items:      items,
		TotalPrice: order.TotalPrice,
		Status:     pb.OrderStatus(order.Status),
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}
}
