package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Các hằng số cấu hình Kafka
const (
	kafkaBroker = "localhost:9092" // Địa chỉ của Kafka Broker
	topic       = "my-first-topic" // Tên topic mà chúng ta sẽ gửi tin nhắn đến
)

func main() {
	// Khởi tạo context
	ctx := context.Background()

	// 1. Tạo một topic nếu nó chưa tồn tại
	// Điều này chỉ là một tiện ích, trong môi trường sản xuất, topic thường được tạo trước
	createTopic(ctx)

	// 2. Tạo một Kafka Writer để gửi tin nhắn
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaBroker),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		BatchTimeout:           5 * time.Millisecond,
		AllowAutoTopicCreation: true, // Cho phép writer tự động tạo topic nếu cần
	}
	defer writer.Close() // Đảm bảo writer được đóng sau khi hoàn thành

	fmt.Printf("Bắt đầu gửi tin nhắn đến topic '%s'...\n", topic)

	// 3. Gửi tin nhắn
	for i := 0; i < 20; i++ {
		// Tạo tin nhắn với Key và Value
		key := fmt.Sprintf("Key-%d", i)
		value := fmt.Sprintf("Tin nhắn số %d từ Go!", i)

		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		}

		err := writer.WriteMessages(ctx, msg)
		if err != nil {
			log.Fatalf("Lỗi khi gửi tin nhắn: %v", err)
		}

		fmt.Printf("Đã gửi tin nhắn: Key='%s', Value='%s'\n", key, value)
		time.Sleep(1 * time.Second) // Đợi 1 giây trước khi gửi tin nhắn tiếp theo
	}

	fmt.Println("Đã gửi tất cả tin nhắn thành công!")
}

// Hàm tiện ích để tạo topic
func createTopic(ctx context.Context) {
	conn, err := kafka.DialContext(ctx, "tcp", kafkaBroker)
	if err != nil {
		log.Fatalf("Lỗi khi kết nối đến Kafka để tạo topic: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("Lỗi khi lấy thông tin controller: %v", err)
	}
	controllerConn, err := kafka.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		log.Fatalf("Lỗi khi kết nối đến controller: %v", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     5,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Printf("Cảnh báo: Không thể tạo topic, có thể đã tồn tại. Lỗi: %v", err)
	} else {
		fmt.Printf("Đã tạo thành công topic '%s'.\n", topic)
	}
}
