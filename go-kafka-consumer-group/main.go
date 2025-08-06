package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

// Các hằng số cấu hình Kafka
const (
	kafkaBroker = "localhost:9092"
	topic       = "test-topic"
	groupID     = "my-consumer-group"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: 10e6,
		// Tắt chế độ tự động commit
		CommitInterval: 0,
		// Quan trọng: Để Commit thủ công, bạn phải tắt CommitInterval.
		// Thư viện sẽ không commit offset nữa, chúng ta phải làm điều đó bằng tay.
	})
	defer reader.Close()

	fmt.Printf("Bắt đầu lắng nghe tin nhắn từ topic '%s' trong Consumer Group '%s' với commit thủ công...\n", topic, groupID)

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		fmt.Println("\nNhận tín hiệu dừng, đang đóng Consumer...")
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer đã đóng.")
			return
		default:
			// Đọc một tin nhắn
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					continue
				}
				log.Fatalf("Lỗi khi đọc tin nhắn: %v", err)
			}

			// Xử lý tin nhắn (giả lập một công việc nặng)
			fmt.Printf("Consumer nhận tin nhắn từ Partition %d, Offset %d, Value=%s - Đang xử lý...\n",
				msg.Partition, msg.Offset, string(msg.Value))
			time.Sleep(2 * time.Second) // Giả lập công việc xử lý mất 2 giây

			// Sau khi xử lý xong, commit offset
			err = reader.CommitMessages(ctx, msg)
			if err != nil {
				log.Fatalf("Lỗi khi commit offset: %v", err)
			}

			fmt.Printf("Đã xử lý xong và commit offset %d\n", msg.Offset)
		}
	}
}
