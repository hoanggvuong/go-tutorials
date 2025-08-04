package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Các hằng số cấu hình Kafka
const (
	kafkaBroker = "localhost:9092"    // Địa chỉ của Kafka Broker
	topic       = "my-first-topic"    // Tên topic mà chúng ta sẽ đọc tin nhắn từ đó
	groupID     = "my-consumer-group" // Tên Consumer Group
)

func main() {
	// 1. Tạo một Kafka Reader để đọc tin nhắn
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker}, // Danh sách các broker để kết nối
		Topic:   topic,                 // Topic để đọc tin nhắn
		GroupID: groupID,               // Tên Consumer Group, giúp Kafka quản lý trạng thái của consumer
		// MinBytes: 10e3, // 10KB
		// MaxBytes: 10e6, // 10MB
		StartOffset: kafka.FirstOffset, // Bắt đầu đọc từ offset đầu tiên của topic
	})
	defer reader.Close() // Đảm bảo reader được đóng sau khi hoàn thành

	fmt.Printf("Bắt đầu lắng nghe tin nhắn từ topic '%s'...\n", topic)

	// 2. Lặp vô hạn để đọc tin nhắn
	for {
		// reader.ReadMessage sẽ block (chờ) cho đến khi có một tin nhắn mới
		// hoặc context bị hủy
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Lỗi khi đọc tin nhắn: %v", err)
		}

		// 3. Xử lý tin nhắn
		// Tin nhắn được đọc sẽ có các thông tin như Topic, Partition, Offset, Key và Value
		fmt.Printf("Nhận tin nhắn: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
}
