package main

import (
	"common/middleware"
	"fmt"
	"log"
	"net"
	"order-service/model"
	"order-service/service"
	"os"
	"strings"

	pb "common/proto/gen/order"
	pbProduct "common/proto/gen/product"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 连接商品服务
	productAddr := os.Getenv("PRODUCT_SERVICE_ADDR")
	productConn, err := grpc.Dial(productAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to product service: %v", err)
	}
	defer productConn.Close()
	productClient := pbProduct.NewProductServiceClient(productConn)

	// 初始化 Redis
	redisAddr := fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	service.InitRedis(redisAddr)

	// 初始化 Kafka
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	// 初始化 Kafka 生产者
	service.InitKafkaProducer(kafkaBrokers, kafkaTopic)
	// 启动 Kafka 消费者
	service.StartKafkaConsumer(kafkaBrokers, kafkaTopic)

	// 创建 gRPC 服务器
	port := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
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
