package service

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

// InitKafkaProducer 初始化 Kafka 生产者
func InitKafkaProducer(brokers []string, topic string) {
	kafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
}

// SendOrderMessage 发送订单消息到 Kafka
func SendOrderMessage(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if kafkaWriter == nil {
		log.Println("Kafka writer is not initialized!")
		return nil
	}
	return kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: []byte(msg),
	})
}
