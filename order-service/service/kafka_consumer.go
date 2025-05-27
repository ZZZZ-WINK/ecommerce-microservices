package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumer(brokers []string, topic string) {
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  "order-service-group",
			MinBytes: 1,
			MaxBytes: 10e6,
		})
		defer r.Close()
		fmt.Println("[Kafka] 消费者已启动，等待订单消息...")
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			m, err := r.ReadMessage(ctx)
			cancel()
			if err != nil {
				if ctx.Err() != context.DeadlineExceeded {
					log.Println("Kafka 消费出错:", err)
				}
				continue
			}
			msg := string(m.Value)
			fmt.Println("[Kafka] 收到消息:", msg)
			// 简单解析消息内容
			if strings.HasPrefix(msg, "order_created") {
				// 这里可以做库存扣减、通知等异步操作
				fmt.Println("[异步处理] 扣减库存 & 发送通知...")
			}
		}
	}()
}
