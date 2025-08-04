package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

// Các hằng số cấu hình Kafka
const (
	kafkaBroker = "localhost:9092"    // Địa chỉ của Kafka Broker
	topic       = "my-first-topic"    // Tên topic mà chúng ta sẽ đọc tin nhắn từ đó
	groupID     = "my-consumer-group" // Tên Consumer Group
)

func main() {
	// Tạo một context để xử lý việc dừng chương trình một cách duyên dáng
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. Tạo một Kafka Reader với GroupID
	// Kafka sẽ tự động phân chia các partitions của topic cho các consumer trong cùng nhóm
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  groupID, // Đây là chìa khóa để tạo Consumer Group
		MaxBytes: 10e6,    // 10MB
	})
	defer reader.Close()

	fmt.Printf("Bắt đầu lắng nghe tin nhắn từ topic '%s' trong Consumer Group '%s'...\n", topic, groupID)

	// Lắng nghe tín hiệu dừng chương trình (như Ctrl+C)
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		fmt.Println("\nNhận tín hiệu dừng, đang đóng Consumer...")
		cancel()
	}()

	// 2. Lặp để đọc tin nhắn cho đến khi context bị hủy
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer đã đóng.")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				// Nếu lỗi là context bị hủy, chúng ta dừng lại
				if ctx.Err() != nil {
					continue
				}
				log.Fatalf("Lỗi khi đọc tin nhắn: %v", err)
			}

			// 3. Xử lý tin nhắn
			fmt.Printf("Consumer Group '%s' nhận tin nhắn từ Partition %d, Offset %d, Key=%s, Value=%s\n",
				groupID, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		}
	}
}
